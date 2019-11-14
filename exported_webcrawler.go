package webcrawler

import (
	"net/http"
)

type Crawler interface {
	Start() error
	Stop()
}

type CrawlerBuilder interface {
	WithMaxDepth(depth int) CrawlerBuilder
	WithFilter(func(crawler Crawler, depth int, url string) bool) CrawlerBuilder
	WithErrorHandler(func(crawler Crawler, err WebCrawlerError)) CrawlerBuilder

	BuildWithOutputDestination(outputDir string) Crawler
	BuildWithOutputHandler(func(crawler Crawler, response *http.Response)) Crawler
}

type WebCrawlerError interface {
	error
}
