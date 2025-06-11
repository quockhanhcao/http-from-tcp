package main

import (
	"log"
	"os"
	"os/signal"
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
	if req.RequestLine.RequestTarget == "/yourproblem" {
		badRequestHandler(w)
	} else if req.RequestLine.RequestTarget == "/myproblem" {
		internalServerErrorHandler(w)
	} else {
		successRequestHandler(w)
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
