package webcrawler

import (
	"testing"
	"net/http"
	"strings"
	"io/ioutil"
	"io"
	"os"
	"path"
	"path/filepath"
)

var indexPageUrl = "http://www.test.com/content.html" 
var indexPageContent = strings.Join([]string {
		"<body>",
		"<a href='/page1.html'>Page 1</a>",
		"<a href='/page2.html'>Page 2</a>",
		"<a href='/subdir/page3.html'>Page 3</a>",
		"</body>",
	},"\n")

func Test_LinksInOtherDomainsRewritten(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "webcrawlerTest")
	if err != nil {
		t.Errorf("Error!! Unable to create temp directory for test. %s", err)
		return
	}
	defer os.Remove(tmpDir)

	bootstrapContent := ".someClass { ... }"
	expectedRequestResponse := map[string]string {
		"http://www.test.com/page1.html": "<link href='https://www.bootstrap.com/bootstrap.min.css' rel='stylessheet'>",
		"https://www.bootstrap.com/bootstrap.min.css": bootstrapContent,
	}
	mockHttpFactory := createMockHttpFactoryWith(expectedRequestResponse)

	builder := createBuilderWith(mockHttpFactory)
	builder.startUrl = "http://www.test.com/page1.html"
	builder.
		BuildWithOutputDestination(tmpDir).
		Start()

	expectedFilePath := path.Join(tmpDir, "www-bootstrap-com/bootstrap.min.css")	
	if content, err := readFileToString(expectedFilePath); err != nil || content != bootstrapContent {
		t.Error("Expected to find bootstrap file")
	}

	expectedContent := "<link href='/www-bootstrap-com/bootstrap.min.css' rel='stylessheet'>" 
	page1Path := path.Join(tmpDir, "page1.html")	
	if page1Content, err := readFileToString(page1Path); err != nil || page1Content != expectedContent {
		t.Error("Expected page1 content to be rewritten")
	}
}

func Test_ErrorHandlerIsCalledOnHttpError(t *testing.T) {
	testWithHttpError(ErrorGet, t)
}

func Test_ErrorHandlerIsCalledOnReadingResponseBody(t *testing.T) {
	testWithHttpError(ErrorRead, t)
}

func Test_ErrorHandlerIsCalledOnBadUrl(t *testing.T) {
	badLink := ":badLink.html"
	badLinks := strings.Join([]string {
		"<body>",
		"<a href='" + badLink  + "'>Bad Link</a>",
		"</body>",
	},"\n")

	expectedRequestResponse := map[string]string {
		indexPageUrl: badLinks,
	}
	mockHttpFactory := createMockHttpFactoryWith(expectedRequestResponse)

	var recordedError WebCrawlerError = nil
	crawlerBuilder := createBuilderWith(mockHttpFactory).
		WithErrorHandler(func(crawler Crawler, err WebCrawlerError) {
			recordedError = err
		}).
		(*crawlerBuilderImpl)
		
	mockOutputHandler := startCrawler(crawlerBuilder)

	test(mockOutputHandler, mockHttpFactory, expectedRequestResponse, t)

	if recordedError == nil {
		t.Errorf("Expected error for %s", badLink)
	}
}

func Test_FilesSavedInCorrectStructure(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "webcrawlerTest")
	if err != nil {
		t.Errorf("Error!! Unable to create temp directory for test. %s", err)
		return
	}

	expectedRequestResponse := map[string]string {
		indexPageUrl: indexPageContent,
		"http://www.test.com/page1.html": "Page 1",
		"http://www.test.com/page2.html": "Page 2",
		"http://www.test.com/subdir/page3.html": "<a href='../page4.html'>Page 4</a>",
		"http://www.test.com/page4.html": "<a href='/page1.html'>Page 1</a>",
	}
	mockHttpFactory := createMockHttpFactoryWith(expectedRequestResponse)

	createBuilderWith(mockHttpFactory).
		BuildWithOutputDestination(tmpDir).
		Start()

	expectedFiles := make([]string, 0, len(expectedRequestResponse))
	for url, expectedContent := range expectedRequestResponse {
		filename := tmpDir + strings.TrimPrefix(url, "http://www.test.com")
		expectedFiles = append(expectedFiles, filename)
		
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			t.Errorf("Couldn't find file %s", filename)
		} else {
			bytes, _ := ioutil.ReadFile(filename)
			strContent := bytesToString(bytes)

			if strContent != expectedContent {
				t.Errorf("In file %s, expecting content:\n%s\n...but got:\n%s", filename, expectedContent, strContent)		
			}
		}
	}

	filepath.Walk(tmpDir, func(path string, info os.FileInfo, err error) error {
		filename := tmpDir + strings.TrimPrefix(path, tmpDir)
		if !info.IsDir() && !contains(expectedFiles, filename) {
			t.Errorf("Found unexpected file: %s", filename)
		}
		return nil
	})

	os.Remove(tmpDir)	
}

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
		BuildWithOutputHandler(func(crawler Crawler, url string, content io.Reader) {
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
		urlsToHttpResponses[k] = NewMockResponse(v, MimeByFilename(k))
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

func testWithHttpError(httpError func(targetUrl string) (*http.Response, error), t *testing.T) {
	var recordedError WebCrawlerError
	crawlerBuilder, _ := NewCrawlerBuilder(indexPageUrl).
		WithErrorHandler(func(crawler Crawler, err WebCrawlerError) {
			recordedError = err
		}).
		(*crawlerBuilderImpl)
	crawlerBuilder.requestFactory = httpError
	
	startCrawler(crawlerBuilder)

	if recordedError == nil {
		t.Error("Expected to received error")
	}
}

func readFileToString(path string) (string, error) {
	if content, err := ioutil.ReadFile(path); err == nil {
		return bytesToString(content), nil
	} else {
		return "", err	
	}
}

func bytesToString(bytes []byte) string {
	strBuilder := &strings.Builder{}
	strBuilder.Write(bytes)
	return strBuilder.String()	
}

func contains(strArr []string, find string) bool {
	for _, v := range strArr {
		if v == find {
			return true
		}
	}
	return false
}

func AssertFileWithContentExists(filename string, expectedContent string, t *testing.T) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("Couldn't find file %s", filename)
	} else {
		bytes, _ := ioutil.ReadFile(filename)
		strContent := bytesToString(bytes)

		if strContent != expectedContent {
			t.Errorf("In file %s, expecting content:\n%s\n...but got:\n%s", filename, expectedContent, strContent)		
		}
	}
}