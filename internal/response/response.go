package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/quockhanhcao/http-from-tcp/internal/headers"
)

type StatusCode int

const (
	StatusOK                StatusCode = 200
	StatusBadRequest        StatusCode = 400
	InternalServerErrorCode StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := ""
	switch statusCode {
	case StatusOK:
		statusLine = "OK"
	case StatusBadRequest:
		statusLine = "Bad Request"
	case InternalServerErrorCode:
		statusLine = "Internal Server Error"
	default:
	}
	w.Write(fmt.Appendf(nil, "HTTP/1.1 %d %s\r\n", statusCode, statusLine))
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers.Set("Content-Length", strconv.Itoa(contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")
	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		_, err := w.Write(fmt.Appendf(nil, "%s: %s\r\n", key, value))
		if err != nil {
			return fmt.Errorf("error writing header %s: %w", key, err)
		}
	}
	w.Write([]byte("\r\n"))
	return nil
}
