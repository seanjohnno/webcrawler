package webcrawler

import (
	"io"
	"io/ioutil"
	"path"
	"os"
	"net/url"
)

type fileOutputHandler struct {
	outputDestination string
	errorHandler func(crawler Crawler, err WebCrawlerError)
}

func (self *fileOutputHandler) ResultHandler(crawler Crawler, rscUrl string, content io.Reader) {
	if parsedUrl, err := url.Parse(rscUrl); err == nil {
		writePath := self.outputDestination + parsedUrl.Path
		
		parentDir := path.Dir(writePath)
		err := os.MkdirAll(parentDir, os.ModePerm)	
		if err != nil {
			// Test	
		}

		byteContent, err := ioutil.ReadAll(content)
		if err != nil {
			// Test
		}
		
		ioutil.WriteFile(writePath, byteContent, 0660)
	} else {
		crawlerErr := createUrlParsingError("", rscUrl, err)
		self.errorHandler(crawler, crawlerErr)
	}	
}