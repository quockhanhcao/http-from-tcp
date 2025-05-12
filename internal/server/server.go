package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/quockhanhcao/http-from-tcp/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
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
	if s.closed.Load() {
		return
	}
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("error accepting connection: %w", err)
			continue
		}
		go s.handle(conn)
	}
}

// handle a single connection by writing the response and closing the connection
func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	// handle the connection
	headers := response.GetDefaultHeaders(0)
	err := response.WriteStatusLine(conn, response.StatusOK)
    if err != nil {
        fmt.Println("error writing status line: %w", err)
    }
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		fmt.Println("error writing headers: %w", err)
	}
}

// Create a net.Listener and returns new Server instance
func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{
		listener: listener,
	}
	server.closed.Store(false)

	go server.listen()
	return server, nil
}
