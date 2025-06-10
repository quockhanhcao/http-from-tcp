package server

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync/atomic"

	"github.com/quockhanhcao/http-from-tcp/internal/request"
	"github.com/quockhanhcao/http-from-tcp/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func (he HandlerError) WriteErrorHandler(w io.Writer) {
	response.WriteStatusLine(w, he.StatusCode)
	messageBytes := []byte(he.Message)
	headers := response.GetDefaultHeaders(len(messageBytes))
	response.WriteHeaders(w, headers)
	w.Write(messageBytes)
}

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

// Close the listener and the server
func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

// uses a loop to .Accept new connections as they come in, handles each one
// in a separate goroutine
func (s *Server) listen(handler Handler) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			fmt.Println("error accepting connection: %w", err)
			continue
		}
		go s.handle(conn, handler)
	}
}

// handle a single connection by writing the response and closing the connection
func (s *Server) handle(conn net.Conn, handler Handler) {
	defer conn.Close()
	// handle the request
	request, err := request.RequestFromReader(conn)
	if err != nil {
		handlerError := &HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    err.Error(),
		}
		handlerError.WriteErrorHandler(conn)
		return
	}
	buf := &bytes.Buffer{}
	handlerError := handler(buf, request)
	if handlerError != nil {
		handlerError.WriteErrorHandler(conn)
		return
	}
	b := buf.Bytes()
	// write the response if succeeded
	response.WriteStatusLine(conn, response.StatusOK)
	headers := response.GetDefaultHeaders(len(b))
	response.WriteHeaders(conn, headers)
	conn.Write(b)
}

// Create a net.Listener and returns new Server instance
func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{
		listener: listener,
		handler:  handler,
	}
	server.closed.Store(false)

	go server.listen(handler)
	return server, nil
}
