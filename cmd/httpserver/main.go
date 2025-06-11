package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/quockhanhcao/http-from-tcp/internal/headers"
	"github.com/quockhanhcao/http-from-tcp/internal/request"
	"github.com/quockhanhcao/http-from-tcp/internal/response"
	"github.com/quockhanhcao/http-from-tcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		proxyHandler(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/yourproblem" {
		badRequestHandler(w)
		return
	} else if req.RequestLine.RequestTarget == "/myproblem" {
		internalServerErrorHandler(w)
		return
	} else {
		successRequestHandler(w)
		return
	}
}

func proxyHandler(w *response.Writer, req *request.Request) {
	requestParams := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	fmt.Println("Proxying request to httpbin with params:", requestParams)
	resp, err := http.Get("https://httpbin.org/" + requestParams)
	if err != nil {
		log.Printf("Error making request to httpbin: %v", err)
		internalServerErrorHandler(w)
		return
	}
	defer resp.Body.Close()

	// for the response to client
	w.WriteStatusLine(response.StatusOK)
	headers := headers.GetDefaultHeaders(0)
	headers.Remove("Content-Length")
	headers.OverrideHeadersByKey("Transfer-Encoding", "chunked")
	w.WriteHeaders(headers)

	buffer := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buffer)
		fmt.Println("Read ", n, "bytes")
		if n > 0 {
			_, err = w.WriteChunkedBody(buffer[:n])
			if err != nil {
				fmt.Printf("Error writing chunked body: %v", err)
				break
			}
		}
		if n == 0 {
			fmt.Println("Reached end of response body from httpbin")
			break
		}
		if err != nil {
			fmt.Printf("Error reading response body from httpbin: %v", err)
			break
		}
	}
	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		fmt.Printf("Error writing chunked body done: %v", err)
	}
}

func badRequestHandler(w *response.Writer) {
	responseStatusCode := response.StatusBadRequest
	responseBody := getResponseBody(response.StatusBadRequest)
	w.WriteStatusLine(responseStatusCode)
	headers := headers.GetDefaultHeaders(len(responseBody))
	headers.OverrideHeadersByKey("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody(responseBody)
}

func internalServerErrorHandler(w *response.Writer) {
	responseStatusCode := response.InternalServerErrorCode
	responseBody := getResponseBody(response.InternalServerErrorCode)
	w.WriteStatusLine(responseStatusCode)
	headers := headers.GetDefaultHeaders(len(responseBody))
	headers.OverrideHeadersByKey("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody(responseBody)
}

func successRequestHandler(w *response.Writer) {
	responseStatusCode := response.StatusOK
	responseBody := getResponseBody(response.StatusOK)
	w.WriteStatusLine(responseStatusCode)
	headers := headers.GetDefaultHeaders(len(responseBody))
	headers.OverrideHeadersByKey("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody(responseBody)
}

func getResponseBody(statusCode response.StatusCode) []byte {
	switch statusCode {
	case response.StatusOK:
		return []byte(`<html>
<head>
    <title>200 OK</title>
</head>
<body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
</body>
</html>`)
	case response.StatusBadRequest:
		return []byte(`<html>
<head>
    <title>400 Bad Request</title>
</head>
<body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
</body>
</html>`)
	case response.InternalServerErrorCode:
		return []byte(`<html>
<head>
    <title>500 Internal Server Error</title>
</head>
<body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
</body>
</html>`)
	default:
		return []byte(`<html>
<head>
    <title>Unknown Status</title>
</head>
<body>
    <h1>Unknown Status</h1>
    <p>Something went wrong.</p>
</body>
</html>`)
	}
}
