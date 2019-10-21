package webcrawler

import (
	"net/http"
	"io/ioutil"
	"regexp"
	"strings"
	"net/url"
)

var startDepth int = 0

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
	
	self.getResource(self.startUrl, startDepth)
}

func (self *webcrawlerImpl) Stop() {
	self.stopped = true
}

func (self *webcrawlerImpl) getResource(url string, depth int) {
	if self.stopped || self.shouldFilter(url) || self.isAlreadyFetched(url) {	
		return
	}
	
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
	
	self.fetchedUrls = append(self.fetchedUrls, url)	
	self.resultHandler(self, url, content)

	nextDepth := depth + 1
	if !self.exceedsMaxDepth(nextDepth) {
		self.recurse(bytesToString(content), url, nextDepth)	
	}
}

func (self *webcrawlerImpl) exceedsMaxDepth(depth int) bool {
	return !(self.maxDepth == -1 || depth  <= self.maxDepth)
} 

func (self *webcrawlerImpl) shouldFilter(url string) bool {
	return self.requestFilter != nil && !self.requestFilter(self, 0, url)
}

func (self *webcrawlerImpl) recurse(content string, currentUrl string, depth int) {
	hrefRegex := regexp.MustCompile("href=['\"](.*)['\"]")
	matches := hrefRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 && len(match[1]) > 0 {
			capturedLink := match[1]
			capturedLink, err := self.getFullUrlPath(capturedLink, currentUrl)
			if err != nil {
				// Test
				return
			}
			self.getResource(capturedLink, depth)
		}
	}	
}

func (self *webcrawlerImpl) getFullUrlPath(capturedLink string, currentUrl string) (string, error) {
	current, err := url.Parse(currentUrl)
	if err != nil {
		return "", err
	}

	toUrl, err := current.Parse(capturedLink)
	if err != nil {
		return "", err
	}

	return toUrl.String(), nil
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
