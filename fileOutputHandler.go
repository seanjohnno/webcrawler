package webcrawler

import (
	"io/ioutil"
	"path"
	"os"
	"net/http"
	"net/url"
	"strings"
	
	"github.com/seanjohnno/webcrawler/linkscanner"
	"github.com/seanjohnno/webcrawler/linkrewriter"
	"github.com/seanjohnno/webcrawler/stringutility"	
)

type fileOutputHandler struct {
	startUrl *url.URL
	outputDestination string
	errorHandler func(crawler Crawler, err WebCrawlerError)

	linkScanner linkscanner.LinkScanner
	linkRewriter linkrewriter.LinkRewriter
}

func (self *fileOutputHandler) ResultHandler(crawler Crawler, response *http.Response) {
	if self.linkScanner == nil {
		self.linkScanner = linkscanner.Create(self.startUrl)
	}

	if self.linkRewriter == nil {
		self.linkRewriter = linkrewriter.Create()
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

	rewriteContext := self.linkRewriter.CreateContext(strContent)
	
	if self.linkScanner.CanScan(response) {
		scanResults := self.linkScanner.Scan(strings.NewReader(strContent), response)
		for _, scanresult := range scanResults {
			if scanresult.Url.Host != self.startUrl.Host {
				rewriteContext.Rewrite(
					scanresult.Url,
					scanresult.Match)
			}	
		}
	}	
	ioutil.WriteFile(writePath, []byte(rewriteContext.String()), 0660)
}

func (self *fileOutputHandler) getWritePath(u *url.URL) string {
	if u.Path == "/" {
		return self.outputDestination + "/index.html" + u.RawQuery
	}

	if rewritten, ok := self.linkRewriter.GetRewrittenUrl(u); ok {
		return self.outputDestination + rewritten
	}

	return self.outputDestination + u.RequestURI()
}
