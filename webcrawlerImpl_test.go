package webcrawler

import (
	"testing"
	"io/ioutil"
	"bytes"
	"net/http"
	"strings"
)

var indexPageUrl = "http://www.test.com/content.html" 
var indexPageContent = strings.Join([]string {
		"<body>",
		"<a href='/page1.html'>Page 1</a>",
		"<a href='/page2.html'>Page 2</a>",
		"<a href='/subdir/page3.html'>Page 3</a>",
		"</body>",
	},"\n")

/*
// fs save
// use mime to search files for correct things i.e. only recurse html/css and css should just be url
// Errors during fetching
// Link rewriting, resources on another server

func Test_FileOutputHandler(t *testing.T) {
	t.Error("Untested")
}

func Test_FontsFetched(t *testing.T) {
	t.Error("Untested")
}
*/

func Test_CssJsFontLinksFetched(t *testing.T)  {
	withJsAndCss := strings.Join([]string {
			"<body>",
			"<script src='/somescript.js'/script>",
			"<link href='/somestyling.css' rel='stylessheet'>",
			"<script src=\"/somescript2.js\"></script>",
			"<link href=\"/somestyling2.css\" rel='stylessheet'>",
			"<a href='/subdir/page3.html'>Page 3</a>",
			"<a href=\"/page4.html\">Page 4</a>",
			"</body>",
		},"\n")
	styling1 := strings.Join([]string {
			"@font-face {",
			"	font-family: Lato;",
			"	src: url(/assets/fonts/Lato.woff2) format('woff2'),",
	        "	 url(/assets/fonts/Lato.woff) format('woff');",
			"}",
		},"\n")
	styling2 := ".header {  background: url(\"/assets/header.png\"); }"  
	expectedRequestResponse := map[string]string {
		indexPageUrl: withJsAndCss,
		"http://www.test.com/somescript.js": "someJavascript() { ... }",
		"http://www.test.com/somestyling.css": styling1,
		"http://www.test.com/somescript2.js": "someJavascript2() { ... }",
		"http://www.test.com/somestyling2.css": styling2,
		"http://www.test.com/subdir/page3.html": "Page 3 content",
		"http://www.test.com/page4.html": "Page 4 content",
		"http://www.test.com/assets/fonts/Lato.woff2": "woff woff 2!",
		"http://www.test.com/assets/fonts/Lato.woff": "woff woff!",
		"http://www.test.com/assets/header.png": "some image content",
	}
	mockHttpFactory := createMockHttpFactoryWith(expectedRequestResponse)
	crawlerBuilder := createBuilderWith(mockHttpFactory)

	mockOutputHandler := startCrawler(crawlerBuilder)
		
	test(mockOutputHandler, mockHttpFactory, expectedRequestResponse, t)
}

func Test_Stop(t *testing.T) {
	expectedRequestResponse := map[string]string {
		indexPageUrl: indexPageContent,
	}
	mockHttpFactory := createMockHttpFactoryWith(expectedRequestResponse)
	crawlerBuilder := createBuilderWith(mockHttpFactory)
	
	requestCount := 0	
	crawler := crawlerBuilder.
		BuildWithOutputHandler(func(crawler Crawler, url string, content []byte) {
			requestCount++
			crawler.Stop()	
		})
	crawler.Start()

	if requestCount > 1 {
		t.Error("Should have only entered handler once")
	}
}

func Test_MaxDepth(t *testing.T) {
	expectedRequestResponse := map[string]string {
		indexPageUrl: indexPageContent,
		"http://www.test.com/page1.html": "Page 1",
		"http://www.test.com/page2.html": "Page 2",
		"http://www.test.com/subdir/page3.html": "<a href='../page4.html'>Page 4</a>",
	}
	mockHttpFactory := createMockHttpFactoryWith(expectedRequestResponse)
	crawlerBuilder := createBuilderWith(mockHttpFactory)
	crawlerBuilder.WithMaxDepth(1)
	
	mockOutputHandler := startCrawler(crawlerBuilder)
	
	test(mockOutputHandler, mockHttpFactory, expectedRequestResponse, t)
}

func Test_FilteringUrls(t *testing.T) {
	expectedRequestResponse := map[string]string {
		indexPageUrl: indexPageContent,
	}
	mockHttpFactory := createMockHttpFactoryWith(expectedRequestResponse)
	crawlerBuilder := createBuilderWith(mockHttpFactory)
	crawlerBuilder.WithFilter(func(crawler Crawler, depth int, path string) bool {
		return !strings.Contains(path, "page")		
	})	

	mockOutputHandler := startCrawler(crawlerBuilder)
	
	test(mockOutputHandler, mockHttpFactory, expectedRequestResponse, t)
}

func Test_LinksAreFetched_AndOnlyOnce(t *testing.T) {
	expectedRequestResponse := map[string]string {
		indexPageUrl: indexPageContent,
		"http://www.test.com/page1.html": "Page 1",
		"http://www.test.com/page2.html": "Page 2",
		"http://www.test.com/subdir/page3.html": "<a href='../page4.html'>Page 4</a>",
		"http://www.test.com/page4.html": "<a href='/page1.html'>Page 1</a>",
	}
	mockHttpFactory := createMockHttpFactoryWith(expectedRequestResponse)
	crawlerBuilder := createBuilderWith(mockHttpFactory)	

	mockOutputHandler := startCrawler(crawlerBuilder)

	test(mockOutputHandler, mockHttpFactory, expectedRequestResponse, t)	
}

func createMockHttpFactoryWith(expectedRequestResponse map[string]string) *mockHttpFactory {
	urlsToHttpResponses := make(map[string]*http.Response)
	for k, v := range expectedRequestResponse {
		urlsToHttpResponses[k] = createHttpResponse(v)
	}
			
	return &mockHttpFactory {
		urlsToResponses: urlsToHttpResponses,
	}
}

func createBuilderWith(httpFactory *mockHttpFactory) *crawlerBuilderImpl {			
	crawlerBuilder, _ := NewCrawlerBuilder(indexPageUrl).(*crawlerBuilderImpl)
	crawlerBuilder.requestFactory = httpFactory.Get
	return crawlerBuilder
}

func startCrawler(crawlerBuilder *crawlerBuilderImpl) *mockOutputHandler {
	mockOutputHandler := &mockOutputHandler {}
	crawler := crawlerBuilder.
		BuildWithOutputHandler(mockOutputHandler.HandleOutput)
	crawler.Start()	

	return mockOutputHandler
}

func test(mockOutputHandler *mockOutputHandler, httpFactory *mockHttpFactory, urlToResponseMap map[string]string, t *testing.T) {	
	for url, expectedContent := range urlToResponseMap {
		requestCount := httpFactory.GetUrlCallCount(url)
		if requestCount != 1 {
			t.Errorf("Expected request to %s once, but was %d times", url, requestCount)
		}
	
		outputContent := mockOutputHandler.GetContentFor(url)
		if outputContent != expectedContent  {
			t.Errorf("Handler. Expected %s but got %s", expectedContent, outputContent)
		}
	}
	
	requestedUrls := httpFactory.requestedUrls
	for _, url := range requestedUrls {
		if urlToResponseMap[url] == "" {
			t.Errorf("Unexpected request to %s", url)
		}
	}
}

func createHttpResponse(response string) *http.Response {
	return &http.Response {
			StatusCode: 200,
			Body: ioutil.NopCloser(
			bytes.NewReader([]byte(response))),
	}
}
