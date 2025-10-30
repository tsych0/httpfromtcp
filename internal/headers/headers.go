package headers

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}

	if idx == 0 {
		return 2, true, nil
	}

	line := string(data[:idx])
	colIdx := strings.Index(line, ":")
	key, val := strings.TrimLeft(line[:colIdx], " "), strings.TrimSpace(line[colIdx+1:])
	if m, _ := regexp.MatchString("^[-!#$%&'*+.^_`|~\\w]*$", key); !m {
		return 0, false, errors.New("error: inalid key")
	}

	key = strings.ToLower(key)
	if v, ok := h[key]; ok {
		h[key] = v + ", " + val
	} else {
		h[key] = val
	}

	return idx + 2, false, nil
}

func (h Headers) Set(key, val string) {
	h[strings.ToLower(key)] = val
}

func (h Headers) Add(key, val string) {
	key = strings.ToLower(key)
	if v, ok := h[key]; ok {
		h[key] = v + ", " + val
	} else {
		h[key] = val
	}
}

func NewHeaders() Headers {
	return make(Headers)
}
