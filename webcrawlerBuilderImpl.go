package webcrawler

import (
	"net/http"
	"net/url"
)

func NewCrawlerBuilder(startUrl string) (CrawlerBuilder, error) {
	if parsedUrl, err := url.Parse(startUrl); err != nil {
		return nil, err
	} else {
		return &crawlerBuilderImpl { 
			startUrl: parsedUrl,
			requestFactory: http.Get,
			maxDepth: -1,
		}, nil
	}
}

type crawlerBuilderImpl struct {
	startUrl *url.URL
	requestFactory func(target string) (*http.Response, error)
	requestFilter func(crawler Crawler, depth int, url *url.URL) bool
	errorHandler func(crawler Crawler, err WebCrawlerError)
	resultHandler func(crawler Crawler, response *http.Response)
	maxDepth int
}

func (self *crawlerBuilderImpl) WithMaxDepth(depth int) CrawlerBuilder {
	self.maxDepth = depth
	return self	
}

func (self *crawlerBuilderImpl) WithFilter(filter func(crawler Crawler, depth int, url *url.URL) bool) CrawlerBuilder {
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
		startUrl: self.startUrl,
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