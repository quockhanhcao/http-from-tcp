package server

import (
	"fmt"
	"github.com/quockhanhcao/http-from-tcp/internal/headers"
	"github.com/quockhanhcao/http-from-tcp/internal/request"
	"github.com/quockhanhcao/http-from-tcp/internal/response"
	"net"
	"sync/atomic"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w *response.Writer, req *request.Request)

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
func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			fmt.Println("error accepting connection: %w", err)
			continue
		}
		go s.handle(conn)
	}
}

// handle a single connection by writing the response and closing the connection
func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	w := response.NewWriter(conn)
	request, err := request.RequestFromReader(conn)
	if err != nil {
		w.WriteStatusLine(response.StatusBadRequest)
		body := []byte(fmt.Sprintf("Error parsing request: %v", err))
		w.WriteHeaders(headers.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
		return
	}
	s.handler(w, request)
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

	go server.listen()
	return server, nil
}
