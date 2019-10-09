package webcrawler

import (
	"net/http"
)

type mockRequest struct {
	requestedUrl string
	mockResponse *http.Response
}

func (self *mockRequest) Get(targetUrl string) (*http.Response, error) {
	self.requestedUrl = targetUrl
	return self.mockResponse, nil
}

type mockOutputHandler struct {
	outputUrl string
	outputBytes []byte
}

func (self *mockOutputHandler) HandleOutput(crawler Crawler, url string, contents []byte) {
	self.outputUrl = url
	self.outputBytes = contents
}
