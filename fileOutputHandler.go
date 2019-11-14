package webcrawler

import (
	"io/ioutil"
	"path"
	"os"
	"net/http"
)

type fileOutputHandler struct {
	outputDestination string
	errorHandler func(crawler Crawler, err WebCrawlerError)
}

func (self *fileOutputHandler) ResultHandler(crawler Crawler, response *http.Response) {
	rscUrl := response.Request.URL.Path
	writePath := self.outputDestination + rscUrl
	
	parentDir := path.Dir(writePath)
	err := os.MkdirAll(parentDir, os.ModePerm)	
	if err != nil {
		self.errorHandler(crawler, createMkDirError(rscUrl, err))
		return
	}

	byteContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		self.errorHandler(crawler, createReadContentError(rscUrl, err))
		return
	}
	
	ioutil.WriteFile(writePath, byteContent, 0660)
}
