package main

import (
	"fmt"
	"log"
	"net"

	"github.com/tsych0/httpfromtcp/internal/request"
)

func main() {
	listner, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", "error", err)
	}
	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
		}
		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("error", err)
		}
		fmt.Println("Request line:")
		fmt.Println("- Method:", req.RequestLine.Method)
		fmt.Println("- Target:", req.RequestLine.RequestTarget)
		fmt.Println("- Version: 1.1")
		fmt.Println("Headers:")
		for k, v := range req.Headers {
			fmt.Printf("- %v: %v\n", k, v)
		}
		fmt.Println("Body:")
		fmt.Println(string(req.Body))
	}
}
