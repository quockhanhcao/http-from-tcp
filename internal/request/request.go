package request

import (
	"fmt"
	"io"
	"slices"
	"strings"
	// "text/template/parse"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var supportedMethods = []string{
	"GET",
	"POST",
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	str, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading request: %v", err)
	}
	requestLine, err := parseRequestLine(string(str))

	if err != nil {
		return nil, fmt.Errorf("error parsing request line: %s", err)
	}

	return &Request{requestLine}, nil
}

func parseRequestLine(requestLine string) (RequestLine, error) {
	requestParts := strings.Split(requestLine, "\r\n")
	requestLine = requestParts[0]
	requestLineParts := strings.Split(requestLine, " ")
	if len(requestLineParts) != 3 {
		return RequestLine{}, fmt.Errorf("invalid request line: %s", requestLine)
	}
	httpMethod := requestLineParts[0]
	if httpMethod != strings.ToUpper(httpMethod) {
		return RequestLine{}, fmt.Errorf("unsupported HTTP method: %s", httpMethod)
	}
	if !slices.Contains(supportedMethods, httpMethod) {
		return RequestLine{}, fmt.Errorf("unsupported HTTP method: %s", httpMethod)
	}
	httpTarget := requestLineParts[1]
	httpVersion := strings.Split(requestLineParts[2], "/")[1]
	if httpVersion != "1.1" {
		return RequestLine{}, fmt.Errorf("not supported HTTP version: %s", httpVersion)
	}

	return RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: httpTarget,
		Method:        httpMethod,
	}, nil
}
