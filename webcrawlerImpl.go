package webcrawler

import (
	"net/http"
)

type webcrawlerImpl struct {
	startUrl string
	requestFactory func(target string) (*http.Response, error)
	requestFilter func(Crawler, int, string) bool
	resultHandler func(Crawler, string, []byte)
	maxDepth int
}

func (self *webcrawlerImpl) Start() {
	
}

func (self *webcrawlerImpl) Stop() {
	
}