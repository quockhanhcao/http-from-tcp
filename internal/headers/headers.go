package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	crlfIdx := bytes.Index(data, []byte(crlf))
	// if no crlf found, return 0 with no error
	// we need to read more data
	if crlfIdx == -1 {
		return 0, false, nil
	}
	if crlfIdx == 0 {
		return 2, true, nil
	}

	parts := bytes.SplitN(data[:crlfIdx], []byte(":"), 2)
	// host key is the Host: part of request header
	hostKey := strings.TrimLeft(string(parts[0]), " ")
	if hostKey != strings.TrimRight(hostKey, " ") {
		return 0, false, fmt.Errorf("invalid header name %s", hostKey)
	}

    hostKey = strings.TrimSpace(hostKey)
	hostAddress := string(parts[1])
	hostAddress = strings.TrimSpace(hostAddress)
    // set the host key to the host address
	h[hostKey] = hostAddress
	return crlfIdx + 2, false, nil
}
