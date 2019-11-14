package webcrawler

import (
	"io"
)

func convertToReaderCloser(reader io.Reader) io.ReadCloser {
	return &readerNopCloser {
		underlyingReader: reader,
	}	
}

type readerNopCloser struct {
	underlyingReader io.Reader	
}

func (self *readerNopCloser) Read(p []byte) (n int, err error) {
	return self.underlyingReader.Read(p)
}

func (self *readerNopCloser) Close() error {
	return nil
}