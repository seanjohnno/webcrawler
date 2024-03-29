package webcrawler

import (
	"testing"
	"io"
	"io/ioutil"
	"os"
	"net/http"
	"net/url"
)

const TestHost string = "http://www.testhost.com/"

func Test_ErrorPassedToHandler_OnContentReadingError(t *testing.T) {
	reportedError := executeWithUrlAndContent("/index.html", &ErrorThrowingReader{}, t)
	if reportedError == nil {
		t.Error("Should have received a content reading error")
	}
}

func Test_MkDir(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "webcrawlerTest")
	if err != nil {
		t.Error("Couldn't create temp directory for test")
	}
	defer os.Remove(tmpDir)

	os.Chmod(tmpDir, 0444)
	defer os.Chmod(tmpDir, 0777)

	testHostUrl, _ := url.Parse(TestHost)

	var reportedError WebCrawlerError = nil
	outputHandler := &fileOutputHandler {
		outputDestination: tmpDir,
		errorHandler: func(crawler Crawler, err WebCrawlerError) {
			reportedError = err
		},
		startUrl: testHostUrl,
	}

	outputHandler.ResultHandler(nil, createResponse("/subDir/test.html", nil))

	if reportedError == nil {
		t.Error("Expected error but got none")
	}
}

func executeWithUrlAndContent(rscUrl string, content io.Reader, t *testing.T) error {
	tmpDir, err := ioutil.TempDir("", "webcrawlerTest")
	if err != nil {
		t.Error("Couldn't create temp directory for test")
		return nil
	}
	defer os.Remove(tmpDir)

	testHostUrl, _ := url.Parse(TestHost)

	var reportedError WebCrawlerError = nil
	outputHandler := &fileOutputHandler {
		outputDestination: tmpDir,
		errorHandler: func(crawler Crawler, err WebCrawlerError) {
			reportedError = err
		},
		startUrl: testHostUrl,
	}

	outputHandler.ResultHandler(nil, createResponse(rscUrl, content))
	return reportedError
}

func createResponse(rscUrl string, content io.Reader) *http.Response {
	baseUrl, _ := url.Parse(TestHost)
	parsedUrl, _ := baseUrl.Parse(rscUrl)	
	return &http.Response {
		Status: "200 OK",
		Body: convertToReaderCloser(content),
		Request: &http.Request {
			URL: parsedUrl,
		},
	}
}