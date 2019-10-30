package webcrawler

import (
	"net/http"
	"net/url"
	"io/ioutil"
	"path"
	"os"
)

type crawlerBuilderImpl struct {
	startUrl string
	requestFactory func(target string) (*http.Response, error)
	requestFilter func(crawler Crawler, depth int, url string) bool
	errorHandler func(crawler Crawler, err WebCrawlerError)
	resultHandler func(crawler Crawler, url string, content []byte)
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

func (self *crawlerBuilderImpl) BuildWithOutputHandler(handler func(crawler Crawler, url string, content []byte)) Crawler {
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

type fileOutputHandler struct {
	outputDestination string
}

func (self *fileOutputHandler) ResultHandler(crawler Crawler, rscUrl string, content []byte) {
	if parsedUrl, err := url.Parse(rscUrl); err == nil {
		writePath := self.outputDestination + parsedUrl.Path
		
		parentDir := path.Dir(writePath)
		err := os.MkdirAll(parentDir, os.ModePerm)	
		if err != nil {
			// Test	
		}
		
		ioutil.WriteFile(writePath, content, 0660)
	} else {
		// Test
	}	
}