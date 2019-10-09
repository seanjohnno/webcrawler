package webcrawler

import (
	"io"
	"net/http"
)

type fileOutputHandler struct {
	outputDestination string
}

func (self *fileOutputHandler) ResultHandler(Crawler, string, []byte) {
}

type crawlerBuilderImpl struct {
	startUrl string
	requestFactory func(method string, target string, body io.Reader) (*http.Request, error)
	requestFilter func(Crawler, int, string) bool
	resultHandler func(Crawler, string, []byte)
	maxDepth int
}

func (self *crawlerBuilderImpl) WithMaxDepth(depth int) CrawlerBuilder {
	self.maxDepth = depth
	return self	
}

func (self *crawlerBuilderImpl) WithFilter(filter func(Crawler, int, string) bool) CrawlerBuilder {
	self.requestFilter = filter
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
		resultHandler: self.resultHandler,
		maxDepth: self.maxDepth,
	}
}

func (self *crawlerBuilderImpl) BuildWithOutputHandler(handler func(Crawler, string, []byte)) Crawler {
	self.resultHandler = handler

	return &webcrawlerImpl {
		startUrl: self.startUrl,
		requestFactory: self.requestFactory,
		requestFilter: self.requestFilter,
		resultHandler: self.resultHandler,
		maxDepth: self.maxDepth,
	}	
}

