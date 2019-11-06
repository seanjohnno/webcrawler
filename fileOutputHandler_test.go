package webcrawler

import (
	"testing"
	"io/ioutil"
	"os"
)

func Test_ErrorPassedToHandler_OnUrlParsingError(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "webcrawlerTest")
	if err != nil {
		t.Error("Couldn't create temp directory for test")
		return
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
		":ICantBeAUrl",
		 nil)

	if reportedError == nil {
		t.Error("Should have received a parsing error")	
	}
}