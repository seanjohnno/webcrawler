package webcrawler

import (
	"testing"
	"io/ioutil"
	"bytes"
	"net/http"
	"strings"
)

func Test_LinksAreFetched_AndOnlyOnce(t *testing.T) {
	startUrl := "http://www.test.com/content.html"
	startContent := strings.Join([]string {
		"<body>",
		"<a href='/page1.html'>Page 1</a>",
		"<a href='/page2.html'>Page 2</a>",
		"<a href='/subdir/page3.html'>Page 3</a>",
		"</body>",
	},"\n")

	expectedRequestResponse := map[string]string {
		startUrl: startContent,
		"http://www.test.com/page1.html": "Page 1",
		"http://www.test.com/page2.html": "Page 2",
		"http://www.test.com/subdir/page3.html": "<a href='../page4.html'>Page 4</a>",
		"http://www.test.com/page4.html": "<a href='/page1.html'>Page 1</a>",
	}
	
	crawler, request, handler := setup(startUrl, expectedRequestResponse)
	crawler.Start()

	for url, expectedContent := range expectedRequestResponse {
		requestCount := request.GetUrlCallCount(url)
		if requestCount != 1 {
			t.Errorf("Expected request to %s once, but was %d times", url, requestCount)
		}

		outputContent := handler.GetContentFor(url)
		if outputContent != expectedContent  {
			t.Errorf("Handler. Expected %s but got %s", expectedContent, outputContent)
		}
	}

	requestedUrls := request.requestedUrls
	for _, url := range requestedUrls {
		if expectedRequestResponse[url] == "" {
			t.Errorf("Unexpected request to %s", url)
		}
	}
}

func setup(startUrl string, urlToResponseMap map[string]string) (Crawler, *mockHttpFactory, *mockOutputHandler) {
	urlsToHttpResponses := make(map[string]*http.Response)
	for k, v := range urlToResponseMap {
		urlsToHttpResponses[k] = createHttpResponse(v)
	}
	
	request := &mockHttpFactory {
		urlsToResponses: urlsToHttpResponses,
	}
	mockOutputHandler := &mockOutputHandler {}
		
	crawlerBuilder := &crawlerBuilderImpl {
		requestFactory: request.Get,
		startUrl: startUrl,
	}
	crawler := crawlerBuilder.
		BuildWithOutputHandler(mockOutputHandler.HandleOutput)
	return crawler, request, mockOutputHandler
}

func createHttpResponse(response string) *http.Response {
	return &http.Response {
			StatusCode: 200,
			Body: ioutil.NopCloser(
			bytes.NewReader([]byte(response))),
	}
}
