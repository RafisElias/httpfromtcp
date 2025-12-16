package request

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

var ErrMalformedRequestLine = errors.New("malformed request line")
var ErrInvalidVersion = errors.New("http version is invalid")
var ErrInvalidMethod = errors.New("method is invalid")

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	requestLine, err := parseRequestLine(data)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *requestLine,
	}, nil
}

func parseRequestLine(data []byte) (*RequestLine, error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return nil, errors.New("not find end of the line")
	}
	currentLine := string(data[:idx])

	stringList := make([]string, 3)

	stringList = strings.Split(currentLine, " ")

	if len(stringList) < 3 {
		return nil, ErrMalformedRequestLine
	}

	method := stringList[0]

	for _, m := range method {
		if m < 'A' || m > 'Z' {
			return nil, ErrInvalidMethod
		}
	}
	requestTarget := stringList[1]
	httpVersion := stringList[2]
	version, _ := strings.CutPrefix(httpVersion, "HTTP/")
	if version != "1.1" {
		return nil, ErrInvalidVersion
	}

	return &RequestLine{
		HttpVersion:   string(version),
		RequestTarget: requestTarget,
		Method:        method,
	}, nil
}
