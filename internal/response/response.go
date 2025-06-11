package response

import (
	"fmt"
	"io"

	"github.com/quockhanhcao/http-from-tcp/internal/headers"
)

type StatusCode int

const (
	StatusOK                StatusCode = 200
	StatusBadRequest        StatusCode = 400
	InternalServerErrorCode StatusCode = 500
)

type WriterState int

const (
	WriterStatusLine WriterState = iota
	WriterHeaders
	WriterBody
)

type Writer struct {
	WriterState WriterState
	writer      io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		WriterState: WriterStatusLine,
		writer:      w,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.WriterState != WriterStatusLine {
		return fmt.Errorf("expected writer state to be writerStatusLine, got %v", w.WriterState)
	}
	_, err := w.writer.Write([]byte(getStatusLine(statusCode)))
	if err != nil {
		return fmt.Errorf("error writing status line: %w", err)
	}
	w.WriterState = WriterHeaders
	return nil
}

func getStatusLine(statusCode StatusCode) string {
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
	return fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, statusLine)
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.WriterState != WriterHeaders {
		return fmt.Errorf("expected writer state to be writeHeaders, got %v", w.WriterState)
	}
	if len(headers) == 0 {
		return fmt.Errorf("headers cannot be empty")
	}
	for key, value := range headers {
		_, err := w.writer.Write(fmt.Appendf(nil, "%s: %s\r\n", key, value))
		if err != nil {
			return fmt.Errorf("error writing header %s: %w", key, err)
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	if err != nil {
		return fmt.Errorf("error writing end of headers: %w", err)
	}
	w.WriterState = WriterBody
	return nil
}

func (w *Writer) WriteBody(body []byte) error {
	if w.WriterState != WriterBody {
		return fmt.Errorf("expected writer state to be writerBody, got %v", w.WriterState)
	}
	if len(body) == 0 {
		return nil
	}
	_, err := w.writer.Write(body)
	if err != nil {
		return fmt.Errorf("error writing body: %w", err)
	}
	return nil
}
