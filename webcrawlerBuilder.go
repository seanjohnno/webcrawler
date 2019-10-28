package webcrawler

import (
	"net/http"
	"net/url"
	"io/ioutil"
	"path"
	"os"
)

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

type crawlerBuilderImpl struct {
	startUrl string
	requestFactory func(target string) (*http.Response, error)
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

