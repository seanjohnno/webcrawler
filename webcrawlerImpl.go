package webcrawler

import (
	"net/http"
	"io/ioutil"
	"github.com/seanjohnno/webcrawler/linkscanner"
)

const startDepth int = 0

type webcrawlerImpl struct {
	startUrl string
	requestFactory func(target string) (*http.Response, error)
	requestFilter func(crawler Crawler, depth int, url string) bool
	errorHandler func(crawler Crawler, err WebCrawlerError)
	resultHandler func(crawler Crawler, url string, content []byte)
	maxDepth int
	
	fetchedUrls []string
	stopped bool
}

func (self *webcrawlerImpl) Start() {
	if self.fetchedUrls == nil {
		self.fetchedUrls = make([]string, 0)	
	}
	
	self.getResource("", self.startUrl, startDepth)
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
		self.errorHandler(self, CreateHttpError(parentUrl, url, err))
		return
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		self.errorHandler(self, CreateHttpError(parentUrl, url, err))
		return
	}
	
	self.fetchedUrls = append(self.fetchedUrls, url)	
	self.resultHandler(self, url, content)

	nextDepth := depth + 1
	if !self.exceedsMaxDepth(nextDepth) {
		self.recurse(parentUrl, content, response, nextDepth)	
	}
}

func (self *webcrawlerImpl) recurse(parentUrl string, responseBody []byte, response *http.Response, depth int) {
	scanResults := linkscanner.Scan(responseBody, response)
	for _, result :=  range scanResults {
		if result.Error == nil {
			self.getResource(parentUrl, result.Url, depth)
		} else {
			self.errorHandler(self, CreateUrlParsingError(parentUrl, result.Url, result.Error))
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
