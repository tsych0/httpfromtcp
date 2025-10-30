package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/tsych0/httpfromtcp/internal/request"
	"github.com/tsych0/httpfromtcp/internal/response"
)

type Server struct {
	port     int
	isClosed atomic.Bool
	listener net.Listener
	handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	listner, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	s := Server{port: port, isClosed: atomic.Bool{}, listener: listner, handler: handler}
	go s.Listen()
	return &s, nil
}

func (s *Server) Close() error {
	s.isClosed.Store(true)
	return nil
}

func (s *Server) Listen() {
	for {
		if s.isClosed.Load() {
			fmt.Println("Closing server")
			break
		}
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatal("Error accepting connection", err)
		}
		s.Handle(conn)
	}
}

func (s *Server) Handle(conn net.Conn) {
	defer conn.Close()
	writer := response.NewWriter(conn)

	req, err := request.RequestFromReader(conn)
	if err != nil {
		m := fmt.Sprintf("Error whlie parsing request: %v", err)
		h := response.GetDefaultHeaders(len(m))
		writer.WriteStatusLine(response.StatusBadRequest)
		writer.WriteHeaders(h)
		writer.WriteBody([]byte(m))
		return
	}

	s.handler(&writer, req)
}
