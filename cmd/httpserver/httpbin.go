package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/tsych0/httpfromtcp/internal/headers"
	"github.com/tsych0/httpfromtcp/internal/request"
	"github.com/tsych0/httpfromtcp/internal/response"
)

func binIt(w *response.Writer, r *request.Request) {
	target := r.RequestLine.RequestTarget
	finalTarget := strings.TrimPrefix(target, "/httpbin/")
	finalUrl := fmt.Sprintf("https://httpbin.org/%v", finalTarget)
	resp, err := http.Get(finalUrl)
	if err != nil {
		myProblem(w, r)
		return
	}
	buf := make([]byte, 1024)
	w.WriteStatusLine(response.StatusCode(resp.StatusCode))

	h := response.GetDefaultHeaders(0)
	h.Add("trailer", "X-Content-SHA256")
	h.Add("trailer", "X-Content-Length")
	h.Set("transfer-encoding", "chunked")
	delete(h, "content-length")
	w.WriteHeaders(h)

	totalBody := []byte{}

	for {
		n, err := resp.Body.Read(buf)
		if err == io.EOF {
			break
		}
		fmt.Printf("Read %v bytes\n", n)
		totalBody = append(totalBody, buf[:n]...)
		w.WriteChunkedBody(buf[:n])
	}
	w.WriteChunkedBodyDone()

	var trailer = make(headers.Headers)
	sha := sha256.Sum256(totalBody)
	trailer.Set("X-Content-SHA256", hex.EncodeToString(sha[:]))
	trailer.Set("X-Content-Length", strconv.Itoa(len(totalBody)))

	w.WriteTrailers(trailer)
}
