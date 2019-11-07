package linkscanner

import (
	"net/url"
	"net/http"
	"regexp"
	"strings"
	"io"
	"io/ioutil"
	)

type ScanResult struct {
	Url string
	Error error
}

type LinkScanner interface {
	CanScan(response *http.Response) bool
	Scan(responseBody io.Reader, response *http.Response) []*ScanResult
}

func Create(origin *url.URL) LinkScanner {
	return &linkScannerImpl {
		originDomain: origin,
	}
}

type linkScannerImpl struct {
	originDomain *url.URL
}

func (_ *linkScannerImpl) CanScan(response *http.Response) bool {
	mimeType := strings.ToLower(mimeFromResponse(response))
	return mimeType == "text/html" || mimeType == "text/css"
}

func (_ *linkScannerImpl) Scan(responseBody io.Reader, response *http.Response) []*ScanResult {
	strResponseBody := responseToString(responseBody)
	mimeType := strings.ToLower(mimeFromResponse(response))

	var regexps []*regexp.Regexp
	if mimeType == "text/html" {
		regexps = []*regexp.Regexp {
			regexp.MustCompile("<a.*?href=['\"](.*?)['\"].*?>"),
			regexp.MustCompile("<link.*?href=['\"](.*?)['\"].*?>"),
			regexp.MustCompile("<script.*?src=['\"](.*?)['\"].*?>"),
			regexp.MustCompile("url\\(['\"]?(.*?)['\"]?\\)"),
		}
	} else if mimeType == "text/css" {
		regexps = []*regexp.Regexp {
			regexp.MustCompile("url\\(['\"]?(.*?)['\"]?\\)"),
		}
	}

	urls := make([]*ScanResult, 0)	
	for _, regexp := range regexps {
		capturedUrls := scanWith(regexp, strResponseBody, response.Request.URL)
		urls = append(urls, capturedUrls...)
	}
	return urls
}

func scanWith(regex *regexp.Regexp, content string, currentUrl *url.URL) []*ScanResult {
	urls := make([]*ScanResult, 0)
	
	matches := regex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 && len(match[1]) > 0 {
			capturedLink := match[1]
			if combinedUrl, err := currentUrl.Parse(capturedLink); err == nil {
				urls = append(urls, &ScanResult {
					Url: combinedUrl.String(),
					Error: err,
				})				
			} else {
				urls = append(urls, &ScanResult {
					Url: capturedLink,
					Error: err,
				})
			}				
		}
	}
	return urls	
}

func mimeFromResponse(response *http.Response) string {
	contentType := response.Header["contentType"]
	if len(contentType) == 0 {
		return ""
	} 
	return contentType[0]
}

func responseToString(content io.Reader) string {
	byteContent, _ := ioutil.ReadAll(content)
	
	strBuilder := &strings.Builder {}
	strBuilder.Write(byteContent)
	return strBuilder.String()
}