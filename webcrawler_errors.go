package webcrawler

import (
	"fmt"
	"net/url"
)

func createUrlParsingError(parentUrl *url.URL, rscUrl *url.URL, innerErr error) WebCrawlerError {
	return &urlParsingError {
		SrcUrl: parentUrl,
		BadUrl: rscUrl, 
		InnerError: innerErr,
	}
}

type urlParsingError struct {
	SrcUrl *url.URL
	BadUrl *url.URL
	InnerError error
}

func (self *urlParsingError) Error() string {
	return fmt.Sprintf("Unable to parse url [%s] at source url [%s]", self.BadUrl, self.SrcUrl)
}

func createHttpError(parentUrl *url.URL, rscUrl *url.URL, innerErr error) WebCrawlerError {
	return &httpGetError {
		SrcUrl: parentUrl,
		ErrorFetchingUrl: rscUrl, 
		InnerError: innerErr,
	}
}

type httpGetError struct {
	SrcUrl *url.URL
	ErrorFetchingUrl *url.URL 
	InnerError error	
}

func (self *httpGetError) Error() string {
	return fmt.Sprintf("Unable to fetch content at [%s]", self.ErrorFetchingUrl)
}

func createReadContentError(url string, innerError error) *ReadContentError {
	return &ReadContentError{
		Url: url,
		InnerError: innerError,
	}
}

type ReadContentError struct {
	Url string
	InnerError error
} 

func (self *ReadContentError) Error() string {
	return fmt.Sprintf("Unable to fetch content at [%s]. Inner error: [%v]", self.Url, self.InnerError)
}

func createMkDirError(url string, innerError error) *MkDirError {
	return &MkDirError{
		Url: url,
		InnerError: innerError,
	}
}

type MkDirError struct {
	Url string
	InnerError error
} 

func (self *MkDirError) Error() string {
	return fmt.Sprintf("Unable to fetch content at [%s]. Inner error: [%v]", self.Url, self.InnerError)
}
