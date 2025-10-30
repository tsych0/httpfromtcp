package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/tsych0/httpfromtcp/internal/request"
	"github.com/tsych0/httpfromtcp/internal/response"
	"github.com/tsych0/httpfromtcp/internal/server"
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

func handler(w *response.Writer, r *request.Request) {
	target := r.RequestLine.RequestTarget
	switch {
	case target == "/yourproblem":
		yourProblem(w, r)
	case target == "/myproblem":
		myProblem(w, r)
	case target == "/video":
		giveTheVideo(w, r)
	case target == "/":
		banger(w, r)
	case strings.HasPrefix(target, "/httpbin/"):
		binIt(w, r)
	default:
		yourProblem(w, r)
	}
}
