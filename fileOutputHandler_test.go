package webcrawler

import (
	"testing"
	"io"
	"io/ioutil"
	"os"
)

func Test_ErrorPassedToHandler_OnUrlParsingError(t *testing.T) {
	reportedError := executeWithUrlAndContent(":ICantBeAUrl", nil, t)
	if reportedError == nil {
		t.Error("Should have received a parsing error")
	}
}

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

	var reportedError WebCrawlerError = nil
	outputHandler := &fileOutputHandler {
		outputDestination: tmpDir,
		errorHandler: func(crawler Crawler, err WebCrawlerError) {
			reportedError = err
		},
	}

	outputHandler.ResultHandler(nil, 
		"/subDir/test.html",
		nil)

	if reportedError == nil {
		t.Error("Expected error but got none")
	}
}

func executeWithUrlAndContent(url string, content io.Reader, t *testing.T) error {
	tmpDir, err := ioutil.TempDir("", "webcrawlerTest")
	if err != nil {
		t.Error("Couldn't create temp directory for test")
		return nil
	}
	defer os.Remove(tmpDir)

	var reportedError WebCrawlerError = nil
	outputHandler := &fileOutputHandler {
		outputDestination: tmpDir,
		errorHandler: func(crawler Crawler, err WebCrawlerError) {
			reportedError = err
		},
	}

	outputHandler.ResultHandler(nil, 
		url,
		content)

	return reportedError
}