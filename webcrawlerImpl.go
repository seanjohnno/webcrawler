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
	startUrl string
	requestFactory func(target string) (*http.Response, error)
	requestFilter func(crawler Crawler, depth int, url string) bool
	errorHandler func(crawler Crawler, err WebCrawlerError)
	resultHandler func(crawler Crawler, response *http.Response)
	maxDepth int

	linkScanner linkscanner.LinkScanner
	
	fetchedUrls []string
	stopped bool
}

func (self *webcrawlerImpl) Start() error {
	if self.fetchedUrls == nil {
		self.fetchedUrls = make([]string, 0)	
	}

	if parsedUrl, err := url.Parse(self.startUrl); err == nil {
		self.linkScanner = linkscanner.Create(parsedUrl)
	} else {
		return err
	}
	
	self.getResource("", self.startUrl, startDepth)
	return nil
}

func (self *webcrawlerImpl) Stop() {
	self.stopped = true
}

func (self *webcrawlerImpl) getResource(parentUrl string, url string, depth int) {
	if self.stopped || self.shouldFilter(url) || self.isAlreadyFetched(url) {	
		return
	}
	
	response, err := self.requestFactory(url)
	if err != nil {
		self.errorHandler(self, createHttpError(parentUrl, url, err))
		return
	}

	if err != nil {
		self.errorHandler(self, createHttpError(parentUrl, url, err))
		return
	}

	self.fetchedUrls = append(self.fetchedUrls, url)	

	nextDepth := depth + 1
	if self.linkScanner.CanScan(response) && !self.exceedsMaxDepth(nextDepth) {
		byteContent, err := ioutil.ReadAll(response.Body)
		if err != nil {
			self.errorHandler(self, createHttpError(parentUrl, url, err))		
			return
		}

		response.Body = convertToReaderCloser(bytes.NewReader(byteContent))
		self.resultHandler(self, response)
		self.recurse(parentUrl, bytes.NewReader(byteContent), response, nextDepth)			
	} else {
		self.resultHandler(self, response)				
	}
}

func (self *webcrawlerImpl) recurse(parentUrl string, responseBody io.Reader, response *http.Response, depth int) {
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

func (self *webcrawlerImpl) shouldFilter(url string) bool {
	return self.requestFilter != nil && !self.requestFilter(self, 0, url)
}

func (self *webcrawlerImpl) isAlreadyFetched(url string) bool {
	for _, v := range self.fetchedUrls {
		if v == url {
			return true
		}
	}	
	return false
}
