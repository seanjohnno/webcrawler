package webcrawler

import (
	"net/http"
	"net/url"
)

type Crawler interface {
	Start()
	Stop()
}

type CrawlerBuilder interface {
	WithMaxDepth(depth int) CrawlerBuilder
	WithFilter(func(crawler Crawler, depth int, url *url.URL) bool) CrawlerBuilder
	WithErrorHandler(func(crawler Crawler, err WebCrawlerError)) CrawlerBuilder

	BuildWithOutputDestination(outputDir string) Crawler
	BuildWithOutputHandler(func(crawler Crawler, response *http.Response)) Crawler
}

type WebCrawlerError interface {
	error
}
