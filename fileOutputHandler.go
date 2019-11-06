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
			self.errorHandler(crawler, createMkDirError(rscUrl, err))
			return
		}

		byteContent, err := ioutil.ReadAll(content)
		if err != nil {
			self.errorHandler(crawler, createReadContentError(rscUrl, err))
			return
		}
		
		ioutil.WriteFile(writePath, byteContent, 0660)
	} else {
		// Don't believe we can ever get here because the url parsing will have failed further up but...
		crawlerErr := createUrlParsingError("", rscUrl, err)
		self.errorHandler(crawler, crawlerErr)
	}	
}
