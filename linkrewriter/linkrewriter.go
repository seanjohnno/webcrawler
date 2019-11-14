package linkrewriter

import (
	"net/url"
	"net/http"
	"io"
)

type Resource struct {
	Url *url.Url
	SavePath string
	Error error
}

type Response struct {
	ModifiedResponseBody io.Reader
	Resources []Resource
}

type LinkRewriter interface {
	Parse(responseBody io.Reader, response *http.Response) (Response, error)
}

func Create(baseUrl *url.URL) LinkRewriter {
	return &linkRewriterImpl {
		BaseUrl: baseUrl,
		Scanner: createScanner(baseUrl),
	}
}

type linkRewriterImpl struct {
	BaseUrl *url.URL	
	Scanner LinkScanner
}

func (self *linkRewriterImpl) Parse(responseBody io.Reader, response *http.Response) (Response, error) {
	if self.Scanner.CanScan(response) {
		if bodyAsString, err := responseToString(responseBody); err != nil {
			return nil, err
		} else {
			resources := make([]Resource, 0)
			scanResults := self.Scanner.Scan(bodyAsString, response)
			for _, scanResult : range scanResults {
				if scanResult.Error != nil {
					resources = append(resources, Resource {
						Url: scanResult.Url,
						Error: Error,
					})			
				} else {
					if BaseUrl.Host != scanResult.Url {
						// Url rewrite
						/*
										Url *url.Url
										SavePath string
										Error error
									
									type ScanResult struct {
										Match string
										Url *url.URL
										Error error
									}
									*/
					} else {
						resources = append(resources, Resource {
							Url: scanResult.Url,
							SavePath: Match,
						})
					}
				}				
			}
			return &Response {
				ModifiedResponseBody: strings.NewReader(bodyAsString),
				Resources: []Resource { }				
			}
		}				
	}
	return &Response {
		ModifiedResponseBody: responseBody,
		Resources: []Resource { }				
	}
}

func ResponseFromScanResult(responseBody string, scanResult ScanResult) (string, *Response) {
	if scanResult.Error != nil {
		return &Resource {
			Url: scanResult.Url,
			SavePath: scanResult.Match,
			Error: scanResult.Error,
		}
	} else if scanResult.Hostname != self.BaseUrl.Hostname {
		rewrittenMatch := rewrite(scanResult.Match)
		bodyAsString = rewriteLinks(bodyAsString, scanResult.Match, rewrittenMatch)
		return &Resource {
			Url: scanResult.Url,
			SavePath: rewrittenMatch,
		}
	} else {
		return &Resource {
			Url: scanResult.Url,
			SavePath: scanResult.Match,
		}
	}
}

func responseToString(content io.Reader) (string, error) {
	if byteContent, err := ioutil.ReadAll(content); err == nil {
		strBuilder := &strings.Builder {}
		strBuilder.Write(byteContent)
		return strBuilder.String()			
	} else {
		return "", err	
	}
}