package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/tsych0/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOk                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) (err error) {
	switch statusCode {
	case StatusOk:
		_, err = fmt.Fprint(w, "HTTP/1.1 200 OK\r\n")
	case StatusBadRequest:
		_, err = fmt.Fprint(w, "HTTP/1.1 400 Bad Request\r\n")
	case StatusInternalServerError:
		_, err = fmt.Fprint(w, "HTTP/1.1 500 Internal Server Error\r\n")
	}
	return
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.Headers{}
	h["content-length"] = strconv.Itoa(contentLen)
	h["connection"] = "close"
	h["content-type"] = "text/plain"
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		_, err := fmt.Fprintf(w, "%v: %v\r\n", k, v)
		if err != nil {
			return err
		}
	}
	fmt.Fprintf(w, "\r\n")
	return nil
}
