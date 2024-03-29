package webcrawler

import (
	"net/http"
	"net/url"
	"errors"
	"strings"
	"bytes"
	"io/ioutil"
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

	requestUrl, _ := url.Parse(targetUrl)
	response := self.urlsToResponses[targetUrl]
	response.Request = &http.Request {
		URL: requestUrl,
	}

	return response, nil
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

func NewMockResponse(body string, mimeType string) *http.Response {
	return &http.Response {
			StatusCode: 200,
			Body: ioutil.NopCloser(
				bytes.NewReader([]byte(body))),
			Header: http.Header {
				"contentType": []string { mimeType },
			},
	}
}

func MimeByFilename(filename string) string {
	parsedUrl, _ := url.Parse(filename)
	
	if strings.HasSuffix(parsedUrl.Path, ".html") {
		return "text/html"
	} else if strings.HasSuffix(parsedUrl.Path, ".js") {
		return "application/javascript"
	} else if strings.HasSuffix(parsedUrl.Path, ".css") {
		return "text/css"
	} else if strings.HasSuffix(parsedUrl.Path, ".png") {
		return "image/png"
	} else {
		return "unknown"
	}
}

type mockOutputHandler struct {
	outputs map[string]string
}

func (self *mockOutputHandler) HandleOutput(crawler Crawler, response *http.Response) {
	if self.outputs == nil {
		self.outputs = make(map[string]string)
	}

	contentAsBytes, _ := ioutil.ReadAll(response.Body)
	strBuilder := &strings.Builder {}
	strBuilder.Write(contentAsBytes)	
	self.outputs[response.Request.URL.String()] = strBuilder.String()
}

func (self *mockOutputHandler) GetContentFor(url string) string {
	return self.outputs[url]
}

func ErrorGet(targetUrl string) (*http.Response, error) {
	return nil, errors.New("Unable to reach server")
}

func ErrorRead(targetUrl string) (*http.Response, error) {
	return &http.Response {				
		Body: &ErrorThrowingReader{},
		Header: map[string][]string {
			"contentType": []string { "text/html", },		
		},
	}, nil
}

type ErrorThrowingReader struct {
}

func (*ErrorThrowingReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("Connection closed") 		
}

func (*ErrorThrowingReader) Close() error {
	return nil
}
