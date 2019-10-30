package webcrawler

import (
	"net/http"
)

type Crawler interface {
	Start()
	Stop()
}

type CrawlerBuilder interface {
	WithMaxDepth(depth int) CrawlerBuilder
	WithFilter(func(crawler Crawler, depth int, url string) bool) CrawlerBuilder
	WithErrorHandler(func(crawler Crawler, err error, url string)) CrawlerBuilder

	BuildWithOutputDestination(outputDir string) Crawler
	BuildWithOutputHandler(func(crawler Crawler, url string, content []byte)) Crawler
}

func NewCrawlerBuilder(url string) CrawlerBuilder {
	return &crawlerBuilderImpl { 
		startUrl: url,
		requestFactory: http.Get,
		maxDepth: -1,
	}
}
