package webcrawler

import (
	"fmt"
)

func createUrlParsingError(src string, errorUrl string, innerErr error) WebCrawlerError {
	return &urlParsingError {
		SrcUrl: src,
		BadUrl: errorUrl, 
		InnerError: innerErr,
	}
}

type urlParsingError struct {
	SrcUrl string
	BadUrl string 
	InnerError error
}

func (self *urlParsingError) Error() string {
	return fmt.Sprintf("Unable to parse url [%s] at source url [%s]", self.BadUrl, self.SrcUrl)
}

func createHttpError(src string, errorUrl string, innerErr error) WebCrawlerError {
	return &httpGetError {
		SrcUrl: src,
		ErrorFetchingUrl: errorUrl, 
		InnerError: innerErr,
	}
}

type httpGetError struct {
	SrcUrl string
	ErrorFetchingUrl string 
	InnerError error	
}

func (self *httpGetError) Error() string {
	return fmt.Sprintf("Unable to fetch content at [%s]", self.ErrorFetchingUrl)
}