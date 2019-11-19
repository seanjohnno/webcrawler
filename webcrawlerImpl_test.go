package webcrawler

import (
	"testing"
	"net/http"
	"net/url"
	"strings"
	"io/ioutil"
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

func Test_RootPathDefaultsToIndexDotHtml(t *testing.T) {
	expectedRequestResponse := map[string]string {
		"http://www.test.com/": "<div>Hey, I'm the index page</div>",
	}

	expectedOutputFiles := map[string]string {
		"/index.html": "<div>Hey, I'm the index page</div>",
	}

	testInputGivesExpectedOutputFiles(
		"http://www.test.com/",
		expectedRequestResponse,
		expectedOutputFiles,
		t)
}

func Test_LinksInOtherDomainsRewritten(t *testing.T) {
	expectedRequestResponse := map[string]string {
		"http://www.test.com/page1.html": "<link href='https://www.bootstrap.com/bootstrap.min.css' rel='stylessheet'>",
		"https://www.bootstrap.com/bootstrap.min.css": ".someClass { ... }",
	}

	expectedOutputFiles := map[string]string {
		"/page1.html": "<link href='/www_bootstrap_com/bootstrap.min.css' rel='stylessheet'>",
		"www_bootstrap_com/bootstrap.min.css": ".someClass { ... }",
	}

	testInputGivesExpectedOutputFiles(
		"http://www.test.com/page1.html",
		expectedRequestResponse,
		expectedOutputFiles,
		t)
}

func Test_QueryStringsWritten(t *testing.T) {
	expectedRequestResponse := map[string]string {
		"http://www.test.com/page1.html?test=one": "<link href='https://www.bootstrap.com/bootstrap.min.css?userId=123' rel='stylessheet'>",
		"https://www.bootstrap.com/bootstrap.min.css?userId=123": ".someClass { ... }",
	}

	expectedOutputFiles := map[string]string {
		"/page1.html?test=one": "<link href='/www_bootstrap_com/bootstrap.min.css?userId=123' rel='stylessheet'>",
		"www_bootstrap_com/bootstrap.min.css?userId=123": ".someClass { ... }",
	}

	testInputGivesExpectedOutputFiles(
		"http://www.test.com/page1.html?test=one",
		expectedRequestResponse,
		expectedOutputFiles,
		t)
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
	expectedRequestResponse := map[string]string {
		indexPageUrl: indexPageContent,
		"http://www.test.com/page1.html": "Page 1",
		"http://www.test.com/page2.html": "Page 2",
		"http://www.test.com/subdir/page3.html": "<a href='../page4.html'>Page 4</a>",
		"http://www.test.com/page4.html": "<a href='/page1.html'>Page 1</a>",
	}

	expectedOutputFiles := map[string]string {
		"/content.html": indexPageContent,
		"/page1.html": "Page 1",
		"/page2.html": "Page 2",
		"/subdir/page3.html": "<a href='../page4.html'>Page 4</a>",
		"/page4.html": "<a href='/page1.html'>Page 1</a>",
	}

	testInputGivesExpectedOutputFiles(
		indexPageUrl,
		expectedRequestResponse,
		expectedOutputFiles,
		t)
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
		BuildWithOutputHandler(func(crawler Crawler, response *http.Response) {
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
	crawlerBuilder.WithFilter(func(crawler Crawler, depth int, path *url.URL) bool {
		return !strings.Contains(path.String(), "page")		
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
	builderInterface, _ := NewCrawlerBuilder(indexPageUrl)
	crawlerBuilder, _ := builderInterface.(*crawlerBuilderImpl)
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
	crawlerBuilderInterface, _ := NewCrawlerBuilder(indexPageUrl)
	crawlerBuilder := crawlerBuilderInterface.
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

func testInputGivesExpectedOutputFiles(startPage string, input map[string]string, outFiles map[string]string, t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "webcrawlerTest")
	if err != nil {
		t.Errorf("Error!! Unable to create temp directory for test. %s", err)
		return
	}
	defer os.Remove(tmpDir)

	mockHttpFactory := createMockHttpFactoryWith(input)

	builder := createBuilderWith(mockHttpFactory)
	parsedUrl, _ := url.Parse(startPage)
	builder.startUrl = parsedUrl
	builder.
		BuildWithOutputDestination(tmpDir).
		Start()

	expectedFiles := make([]string, 0, len(outFiles))
	for filepath, expectedContent := range outFiles {
		filepath = path.Join(tmpDir, filepath)
		expectedFiles = append(expectedFiles, filepath)
		
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			t.Errorf("Couldn't find file %s", filepath)
		} else {
			bytes, _ := ioutil.ReadFile(filepath)
			strContent := bytesToString(bytes)

			if strContent != expectedContent {
				t.Errorf("In file %s, expecting content:\n%s\n...but got:\n%s", filepath, expectedContent, strContent)		
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