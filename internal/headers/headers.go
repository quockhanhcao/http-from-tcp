package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

var alphanumeric = regexp.MustCompile("^[A-Za-z0-9!#$%&'*+\\-.\\^_`|~]+$")

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
	// if we find a CRLF at the beginning of the data, we are done
	if crlfIdx == 0 {
		return 2, true, nil
	}

	parts := bytes.SplitN(data[:crlfIdx], []byte(":"), 2)
	// host key is the "Host:" or "User-Agent" part of request header
	hostKey := strings.TrimLeft(string(parts[0]), " ")
	if hostKey != strings.TrimRight(hostKey, " ") {
		return 0, false, fmt.Errorf("invalid header name %s", hostKey)
	}

	// can use isAlphaNumeric or validString to check valid header
	if !validString(hostKey) || len(hostKey) == 0 {
		return 0, false, fmt.Errorf("invalid header name %s", hostKey)
	}

	hostKey = strings.TrimSpace(hostKey)
	hostKey = strings.ToLower(hostKey)

	// process host value
	hostValues := string(parts[1])
	hostValues = strings.TrimSpace(hostValues)
	h.Set(hostKey, hostValues)

	return crlfIdx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	v, ok := h[key]
	if ok {
		value = v + ", " + value
	}
	h[key] = value
}

func (h Headers) Get(key string) (string, bool) {
    key = strings.ToLower(key)
	v, ok := h[key]
	return v, ok
}

func isAlphaNumeric(s string) bool {
	return alphanumeric.MatchString(s)
}

func validString(s string) bool {
	for _, char := range s {
		if !(char >= 'A' && char <= 'Z' || char >= 'a' && char <= 'z' || char >= '0' && char <= '9' || char == '-') {
			return false
		}
	}
	return true
}
