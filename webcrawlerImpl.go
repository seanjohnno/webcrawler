package webcrawler

import (
	"io"
	"net/http"
)

type webcrawlerImpl struct {
	startUrl string
	requestFactory func(method string, target string, body io.Reader) (*http.Request, error)
	requestFilter func(Crawler, int, string) bool
	resultHandler func(Crawler, string, []byte)
	maxDepth int
}

func (self *webcrawlerImpl) Start() {
	
}

func (self *webcrawlerImpl) Stop() {
	
}