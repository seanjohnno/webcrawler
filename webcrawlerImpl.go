package webcrawler

import (
	"net/http"
	"io/ioutil"
	"regexp"
	"strings"
	"net/url"
)

type webcrawlerImpl struct {
	startUrl string
	requestFactory func(target string) (*http.Response, error)
	requestFilter func(Crawler, int, string) bool
	resultHandler func(Crawler, string, []byte)
	maxDepth int

	fetchedUrls []string
	stopped bool
}

func (self *webcrawlerImpl) Start() {
	if self.fetchedUrls == nil {
		self.fetchedUrls = make([]string, 0)	
	}
	
	self.getResource(self.startUrl, 0)
}

func (self *webcrawlerImpl) Stop() {
	self.stopped = true
}

func (self *webcrawlerImpl) getResource(url string, depth int) {
	if self.stopped == true {
		return
	}

	if self.requestFilter != nil && !self.requestFilter(self, 0, url) {
		return
	}

	if self.isAlreadyFetched(url) {
		return
	}
	self.fetchedUrls = append(self.fetchedUrls, url)

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

	if self.maxDepth == -1 || depth + 1 <= self.maxDepth {
		self.recurse(bytesToString(content), url, depth + 1)	
	}
}

func (self *webcrawlerImpl) recurse(content string, currentUrl string, depth int) {
	hrefRegex := regexp.MustCompile("href=['\"](.*)['\"]")
	matches := hrefRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 && len(match[1]) > 0 {
			capturedLink := match[1]
			capturedLink = self.getFullUrlPath(capturedLink, currentUrl)
			self.getResource(capturedLink, depth)
		}
	}	
}

func (self *webcrawlerImpl) getFullUrlPath(capturedLink string, currentUrl string) string {
	current, err := url.Parse(currentUrl)
	if err != nil {
		return ""
	}

	toUrl, err := current.Parse(capturedLink)
	if err != nil {
		return ""
	}

	return toUrl.String()
}

func (self *webcrawlerImpl) isAlreadyFetched(url string) bool {
	for _, v := range self.fetchedUrls {
		if v == url {
			return true
		}
	}	
	return false
}

func bytesToString(content []byte) string {
	strBuilder := &strings.Builder {}
	strBuilder.Write(content)
	return strBuilder.String()
}
