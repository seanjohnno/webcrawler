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
	WithFilter(filter func(Crawler, int, string) bool) CrawlerBuilder
	BuildWithOutputDestination(string) Crawler
	BuildWithOutputHandler(handler func(Crawler, string, []byte)) Crawler
}

func NewCrawlerBuilder(url string) CrawlerBuilder {
	return &crawlerBuilderImpl { 
		startUrl: url,
		requestFactory: http.Get,
	}
}

