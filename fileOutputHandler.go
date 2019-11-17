package webcrawler

import (
	"io/ioutil"
	"path"
	"os"
	"net/http"
	"net/url"
	"strings"
	
	"github.com/seanjohnno/webcrawler/linkscanner"
	"github.com/seanjohnno/webcrawler/stringutility"	
)

type fileOutputHandler struct {
	startUrl *url.URL
	outputDestination string
	errorHandler func(crawler Crawler, err WebCrawlerError)

	linkScanner linkscanner.LinkScanner
	rewrittenLinks []*url.URL
}

func (self *fileOutputHandler) ResultHandler(crawler Crawler, response *http.Response) {
	if self.linkScanner == nil {
		self.linkScanner = linkscanner.Create(self.startUrl)
	}
	
	if self.rewrittenLinks == nil {
		self.rewrittenLinks = make([]*url.URL, 0)
	}

	writePath := self.getWritePath(response.Request.URL)
					
	parentDir := path.Dir(writePath)
	err := os.MkdirAll(parentDir, os.ModePerm)	
	if err != nil {
		self.errorHandler(crawler, createMkDirError(writePath, err))
		return
	}

	strContent, err := stringutility.ReaderToString(response.Body)
	if err != nil {
		self.errorHandler(crawler, createReadContentError(response.Request.URL.String(), err))
		return
	}

	if self.linkScanner.CanScan(response) {
		scanResults := self.linkScanner.Scan(strings.NewReader(strContent), response)
		for _, scanresult := range scanResults {
			if scanresult.Url.Host != self.startUrl.Host {
				strContent = strings.ReplaceAll(
					strContent,
					 scanresult.Match, 
					 self.rewriteUrl(
					 	scanresult.Url))
				
				self.rewrittenLinks = append(self.rewrittenLinks, scanresult.Url)			
			}	
		}
	}	
	ioutil.WriteFile(writePath, []byte(strContent), 0660)
}

func (self *fileOutputHandler) getWritePath(u *url.URL) string {
	if self.shouldBeRewritten(u) {
		return self.rewriteUrl(u)
	} else {
		return self.outputDestination + u.Path
	}
}

func (self *fileOutputHandler) rewriteUrl(link *url.URL) string {
	return strings.ReplaceAll(link.Host, ".", "_") + link.Path  
}

func (self *fileOutputHandler) shouldBeRewritten(u *url.URL) bool {
	for _, rewriteUrl := range self.rewrittenLinks {
		if rewriteUrl.String() == u.String() {
			return true
		}
	}
	return false
}