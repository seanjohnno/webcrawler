package webcrawler

import (
	"net/http"
	"fmt"
	"io"
)

type Crawler interface {
	Start()
	Stop()
}

type CrawlerBuilder interface {
	WithMaxDepth(depth int) CrawlerBuilder
	WithFilter(func(crawler Crawler, depth int, url string) bool) CrawlerBuilder
	WithErrorHandler(func(crawler Crawler, err WebCrawlerError)) CrawlerBuilder

	BuildWithOutputDestination(outputDir string) Crawler
	BuildWithOutputHandler(func(crawler Crawler, url string, content io.Reader)) Crawler
}

func NewCrawlerBuilder(url string) CrawlerBuilder {
	return &crawlerBuilderImpl { 
		startUrl: url,
		requestFactory: http.Get,
		maxDepth: -1,
	}
}

type WebCrawlerError interface {
	error
}

func CreateUrlParsingError(src string, errorUrl string, innerErr error) WebCrawlerError {
	return &UrlParsingError {
		SrcUrl: src,
		BadUrl: errorUrl, 
		InnerError: innerErr,
	}
}

type UrlParsingError struct {
	SrcUrl string
	BadUrl string 
	InnerError error
}

func (self *UrlParsingError) Error() string {
	return fmt.Sprintf("Unable to parse url [%s] at source url [%s]", self.BadUrl, self.SrcUrl)
}

func CreateHttpError(src string, errorUrl string, innerErr error) WebCrawlerError {
	return &HttpGetError {
		SrcUrl: src,
		ErrorFetchingUrl: errorUrl, 
		InnerError: innerErr,
	}
}

type HttpGetError struct {
	SrcUrl string
	ErrorFetchingUrl string 
	InnerError error	
}

func (self *HttpGetError) Error() string {
	return fmt.Sprintf("Unable to fetch content at [%s]", self.ErrorFetchingUrl)
}