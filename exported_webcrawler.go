package webcrawler

import (
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

type WebCrawlerError interface {
	error
}
