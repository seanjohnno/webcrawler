package stringutility

import (
	"io"
	"io/ioutil"
	"strings"
)


func bytesToString(byteArr []byte) string {
	builder := &strings.Builder{}
	builder.Write(byteArr)
	return builder.String()
}

func ReaderToString(reader io.Reader) (string, error) {
	if bytes, err := ioutil.ReadAll(reader); err != nil {
		return "", err
	} else {
		return bytesToString(bytes), nil
	}
}