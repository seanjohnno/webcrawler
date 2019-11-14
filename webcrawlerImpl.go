package webcrawler

import (
	"net/http"
	"net/url"
	"io/ioutil"
	"io"
	"github.com/seanjohnno/webcrawler/linkscanner"
	"bytes"
)

const startDepth int = 0

type webcrawlerImpl struct {
	startUrl *url.URL
	requestFactory func(target string) (*http.Response, error)
	requestFilter func(crawler Crawler, depth int, url *url.URL) bool
	errorHandler func(crawler Crawler, err WebCrawlerError)
	resultHandler func(crawler Crawler, response *http.Response)
	maxDepth int

	linkScanner linkscanner.LinkScanner
	
	fetchedUrls []*url.URL
	stopped bool
}

func (self *webcrawlerImpl) Start() {
	if self.fetchedUrls == nil {
		self.fetchedUrls = make([]*url.URL, 0)	
	}

	self.linkScanner = linkscanner.Create(self.startUrl)
	
	self.getResource(nil, self.startUrl, startDepth)
}

func (self *webcrawlerImpl) Stop() {
	self.stopped = true
}

func (self *webcrawlerImpl) getResource(parentUrl *url.URL, rscUrl *url.URL, depth int) {
	if self.stopped || self.shouldFilter(rscUrl) || self.isAlreadyFetched(rscUrl) {	
		return
	}
	
	response, err := self.requestFactory(rscUrl.String())
	if err != nil {
		self.errorHandler(self, createHttpError(parentUrl, rscUrl, err))
		return
	}

	if err != nil {
		self.errorHandler(self, createHttpError(parentUrl, rscUrl, err))
		return
	}

	self.fetchedUrls = append(self.fetchedUrls, rscUrl)	

	nextDepth := depth + 1
	if self.linkScanner.CanScan(response) && !self.exceedsMaxDepth(nextDepth) {
		byteContent, err := ioutil.ReadAll(response.Body)
		if err != nil {
			self.errorHandler(self, createHttpError(parentUrl, rscUrl, err))		
			return
		}

		response.Body = convertToReaderCloser(bytes.NewReader(byteContent))
		self.resultHandler(self, response)
		self.recurse(parentUrl, bytes.NewReader(byteContent), response, nextDepth)			
	} else {
		self.resultHandler(self, response)				
	}
}

func (self *webcrawlerImpl) recurse(parentUrl *url.URL, responseBody io.Reader, response *http.Response, depth int) {
	scanResults := self.linkScanner.Scan(responseBody, response)
	for _, result :=  range scanResults {
		if result.Error == nil {
			self.getResource(parentUrl, result.Url, depth)
		} else {
			self.errorHandler(self, createUrlParsingError(parentUrl, result.Url, result.Error))
		}
	}
}

func (self *webcrawlerImpl) exceedsMaxDepth(depth int) bool {
	return !(self.maxDepth == -1 || depth  <= self.maxDepth)
} 

func (self *webcrawlerImpl) shouldFilter(url *url.URL) bool {
	return self.requestFilter != nil && !self.requestFilter(self, 0, url)
}

func (self *webcrawlerImpl) isAlreadyFetched(url *url.URL) bool {
	for _, v := range self.fetchedUrls {
		if v.String() == url.String() {
			return true
		}
	}	
	return false
}
