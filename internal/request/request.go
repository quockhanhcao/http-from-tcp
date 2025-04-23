package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/quockhanhcao/http-from-tcp/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	// track the state of the parser
	// 0: intialized state
	// 1: done
	State requestState
	// Headers
	Headers headers.Headers
}

type requestState int

const (
	requestStateIntialized requestState = iota
	requestStateParsingHeaders
	requestStateDone
)

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
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	// instead of reading all the data into memory
	// read it into chunks
	// and parse it as we go
	buffer := make([]byte, bufferSize, bufferSize)
	// keep track of how much data we have read
	readToIndex := 0
	// new request
	request := &Request{
		State:   requestStateIntialized,
		Headers: headers.NewHeaders(),
	}

	for request.State != requestStateDone {
		if readToIndex >= len(buffer) {
			// the buffer is full
			// grow it
			newBuffer := make([]byte, 2*len(buffer))
			copy(newBuffer, buffer)
			buffer = newBuffer
		}
		// read from the reader to the buffer, from readToIndex index position to add in more data to the buffer
		numBytesRead, err := reader.Read(buffer[readToIndex:])
		if errors.Is(err, io.EOF) {
			if request.State != requestStateDone {
                return nil, fmt.Errorf("request not complete")
			}
		}
		// increase the readToIndex by the number of bytes read
		// later, with new iteration, we can check whether if the buffer has to be resized
		readToIndex += numBytesRead
		numBytesParse, err := request.parse(buffer[:readToIndex])
		if err != nil {
			return nil, err
		}
		// remove the data we have parsed
		// this keep the buffer small and memory efficient
		copy(buffer, buffer[numBytesParse:])
		readToIndex -= numBytesParse
	}
	return request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParse := 0
	for r.State != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParse:])
		if err != nil {
			return 0, err
		}
		totalBytesParse += n
		if n == 0 {
			break
		}
	}
	return totalBytesParse, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.State {
	case requestStateIntialized:
		requestLine, parsedBytes, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if parsedBytes == 0 {
			// no byte is parsed, we need more data
			return 0, nil
		}

		r.RequestLine = *requestLine
		r.State = requestStateParsingHeaders
		return parsedBytes, nil

	case requestStateParsingHeaders:
		// parse headers
		// keep track of the number of bytes parsed
		totalBytesParsed, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.State = requestStateDone
		}
		return totalBytesParsed, nil
	case requestStateDone:
		return 0, fmt.Errorf("trying to read data in a done state")
	default:
		return 0, fmt.Errorf("unknown state")
	}
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}
	// plus 2 for the CRLF
	parsedBytes := idx + 2
	return requestLine, parsedBytes, nil
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
