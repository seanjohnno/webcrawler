package webcrawler

import (
	"net/http"
	"errors"
	"strings"
)

type mockHttpFactory struct {
	urlsToResponses map[string]*http.Response
	requestedUrls []string
}

func (self *mockHttpFactory) Get(targetUrl string) (*http.Response, error) {
	if self.requestedUrls == nil {
		self.requestedUrls = make([]string, 0)
	}

	self.requestedUrls = append(self.requestedUrls, targetUrl)
	return self.urlsToResponses[targetUrl], nil
}

func (self *mockHttpFactory) GetUrlCallCount(targetUrl string) int {
	count := 0
	for _, v := range self.requestedUrls {
		if v == targetUrl {
			count++
		}
	}
	return count
}

type mockOutputHandler struct {
	outputs map[string]string
}

func (self *mockOutputHandler) HandleOutput(crawler Crawler, url string, contents []byte) {
	if self.outputs == nil {
		self.outputs = make(map[string]string)
	}

	strBuilder := &strings.Builder {}
	strBuilder.Write(contents)	
	self.outputs[url] = strBuilder.String()
}

func (self *mockOutputHandler) GetContentFor(url string) string {
	return self.outputs[url]
}

func ErrorGet(targetUrl string) (*http.Response, error) {
	return nil, errors.New("Unable to reach server")
}