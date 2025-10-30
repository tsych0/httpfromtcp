package server

import (
	"github.com/tsych0/httpfromtcp/internal/request"
	"github.com/tsych0/httpfromtcp/internal/response"
)

// type HandlerError struct {
// 	StatusCode response.StatusCode
// 	Message    string
// }

type Handler func(w *response.Writer, req *request.Request)

// func (h *HandlerError) Write(w *response.Writer) {
// 	w.WriteStatusLine(h.StatusCode)
// 	w.WriteHeaders(response.GetDefaultHeaders(len(h.Message)))
// 	w.WriteBody([]byte(h.Message))
// }
