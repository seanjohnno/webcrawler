package webcrawler

import (
	"net/http"
	"io/ioutil"
	"regexp"
	"strings"
)

type webcrawlerImpl struct {
	startUrl string
	requestFactory func(target string) (*http.Response, error)
	requestFilter func(Crawler, int, string) bool
	resultHandler func(Crawler, string, []byte)
	maxDepth int
}

func (self *webcrawlerImpl) Start() {
	self.getResource(self.startUrl)
}

func (self *webcrawlerImpl) Stop() {
	// Test stop before end
}

func (self *webcrawlerImpl) getResource(url string) {
	// Test request filter
	response, err := self.requestFactory(url)
	if err != nil {
		// Test
		return
	}

	if response == nil {
		// Test
		return
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		// Test
		return
	}
	self.resultHandler(self, url, content)
	self.recurse(bytesToString(content), url)
}

func (self *webcrawlerImpl) recurse(content string, currentUrl string) {
	hrefRegex := regexp.MustCompile("href=['\"](.*)['\"]")
	matches := hrefRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 && len(match[1]) > 0 {
			capturedLink := match[1]
			capturedLink = self.getFullUrlPath(capturedLink, currentUrl)
			self.getResource(capturedLink)
		}
	}	
}

func (self *webcrawlerImpl) getFullUrlPath(capturedLink string, currentUrl string) string {
	
	return currentUrl + capturedLink
}

func bytesToString(content []byte) string {
	strBuilder := &strings.Builder {}
	strBuilder.Write(content)
	return strBuilder.String()
}
