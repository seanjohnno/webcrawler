package webcrawler

import (
	"net/http"
)

func NewCrawlerBuilder(url string) CrawlerBuilder {
	return &crawlerBuilderImpl { 
		startUrl: url,
		requestFactory: http.Get,
		maxDepth: -1,
	}
}

type crawlerBuilderImpl struct {
	startUrl string
	requestFactory func(target string) (*http.Response, error)
	requestFilter func(crawler Crawler, depth int, url string) bool
	errorHandler func(crawler Crawler, err WebCrawlerError)
	resultHandler func(crawler Crawler, response *http.Response)
	maxDepth int
}

func (self *crawlerBuilderImpl) WithMaxDepth(depth int) CrawlerBuilder {
	self.maxDepth = depth
	return self	
}

func (self *crawlerBuilderImpl) WithFilter(filter func(crawler Crawler, depth int, url string) bool) CrawlerBuilder {
	self.requestFilter = filter
	return self	
}

func (self *crawlerBuilderImpl) WithErrorHandler(errHandler func(crawler Crawler, err WebCrawlerError)) CrawlerBuilder {
	self.errorHandler = errHandler
	return self
}

func (self *crawlerBuilderImpl) BuildWithOutputDestination(destination string) Crawler {
	fileOutputRequestHandler := &fileOutputHandler {
		outputDestination: destination,
	}
	self.resultHandler = fileOutputRequestHandler.ResultHandler 	

	return &webcrawlerImpl {
		startUrl: self.startUrl,
		requestFactory: self.requestFactory,
		requestFilter: self.requestFilter,
		errorHandler: self.errorHandler,
		resultHandler: self.resultHandler,
		maxDepth: self.maxDepth,
	}
}

func (self *crawlerBuilderImpl) BuildWithOutputHandler(handler func(crawler Crawler, response *http.Response)) Crawler {
	self.resultHandler = handler

	return &webcrawlerImpl {
		startUrl: self.startUrl,
		requestFactory: self.requestFactory,
		requestFilter: self.requestFilter,
		errorHandler: self.errorHandler,
		resultHandler: self.resultHandler,
		maxDepth: self.maxDepth,
	}	
}