package request

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/tsych0/httpfromtcp/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	state       RequestState
}

type RequestState int

const (
	requestStateInitialized RequestState = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const bufferSize = 8

func RequestFromReader(f io.Reader) (*Request, error) {
	r := Request{
		Headers: headers.NewHeaders(),
		state:   requestStateInitialized,
	}
	bytesRead := 0
	bytesParsed := 0
	buf := make([]byte, bufferSize)

	for !(r.state == requestStateDone) {
		if bytesRead >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		n, err := f.Read(buf[bytesRead:])
		if err == io.EOF {
			if r.state != requestStateDone {
				return nil, errors.New("got eof before finishing")
			}
			break
		}

		if err != nil {
			return nil, err
		}

		bytesRead += n
		for {
			parsed, err := r.parse(buf[:bytesRead])
			if err != nil {
				return nil, err
			}
			if parsed == 0 {
				break
			}
			bytesRead -= parsed
			bytesParsed += parsed
			copy(buf, buf[parsed:])
		}

	}

	return &r, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.state == requestStateInitialized {
		rl, consumed, err := parseRequestLine(data)
		if consumed == 0 && err == nil {
			return 0, nil
		}
		r.state = requestStateParsingHeaders
		if err != nil {
			return 0, err
		}
		r.RequestLine = *rl
		return consumed, nil
	}

	if r.state == requestStateParsingHeaders {
		cons, done, err := r.Headers.Parse(data)
		if err != nil {
			r.state = requestStateDone
			return 0, err
		}

		if done {
			r.state = requestStateParsingBody
			return cons, nil
		}

		return cons, nil
	}

	if r.state == requestStateParsingBody {
		if length, ok := r.Headers["content-length"]; ok {
			l, err := strconv.Atoi(length)
			if err != nil {
				return 0, err
			}

			r.Body = append(r.Body, data...)
			if len(r.Body) > l {
				return 0, errors.New("more data given than mentioned")
			}

			if l == len(r.Body) {
				r.state = requestStateDone
			}
			return len(data), nil
		} else {
			r.state = requestStateDone
		}
	}

	return 0, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return nil, 0, nil
	}

	line := string(data[:idx])

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, 0, errors.New("error: invalid request")
	}

	method := parts[0]
	if method != strings.ToUpper(method) {
		return nil, 0, errors.New("error: invalid method")
	}

	target := parts[1]
	version := parts[2]
	if version != "HTTP/1.1" {
		return nil, 0, errors.New("error: not supported")
	}

	httpVersion := strings.Split(version, "/")[1]

	r := RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: target,
		Method:        method,
	}

	return &r, idx + 2, nil
}
