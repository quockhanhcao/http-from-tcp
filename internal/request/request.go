package request

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"strings"
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

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	str, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading request: %v", err)
	}
	requestLine, err := parseRequestLine(str)

	if err != nil {
		return nil, fmt.Errorf("error parsing request line: %s", err)
	}

	return &Request{*requestLine}, nil
}

func parseRequestLine(data []byte) (*RequestLine, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, fmt.Errorf("could not find CRLF in request-line")
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, err
	}
	return requestLine, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("poor format request-line %s", str)
	}
	httpMethod := parts[0]
	if !slices.Contains(supportedMethods, httpMethod) {
		return nil, fmt.Errorf("invalid method %s", httpMethod)
	}

	httpTarget := parts[1]

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line %s", parts[2])
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized http version %s", httpPart)
	}
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized http version %s", version)
	}

	return &RequestLine{
		Method:        httpMethod,
		RequestTarget: httpTarget,
		HttpVersion:   version,
	}, nil
}
