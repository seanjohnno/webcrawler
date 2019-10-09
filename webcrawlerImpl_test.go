package webcrawler

import (
	"testing"
)

func Test_WebRequestIsSentToCorrectStartUrl(t *testing.T) {
		crawler.Start()
}

func setup() (Crawler, *mockRequest, *mockOutputHandler) {
	request := &mockRequest {
		}
		
		mockOutputHandler := &mockOutputHandler {
		}
		
		crawlerBuilder := &crawlerBuilderImpl {
			requestFactory: request.Get,
		}			
		crawler := crawlerBuilder.
			BuildWithOutputHandler(mockOutputHandler.HandleOutput)
	return crawler
}

/*type CrawlerBuilder interface {
	WithMaxDepth(depth int) CrawlerBuilder
	WithFilter(filter func(Crawler, int, string) bool) CrawlerBuilder
	BuildWithOutputDestination(string) Crawler
	BuildWithOutputHandler(handler func(Crawler, string, []byte)) Crawler
}*/

