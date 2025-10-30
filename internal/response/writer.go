package response

import (
	"bytes"
	"fmt"
	"io"

	"github.com/tsych0/httpfromtcp/internal/headers"
)

type WriterState int

const (
	WriterStateInitialized WriterState = iota
	WriterStateWritingHeader
	WriterStateWritingBody
	writerStateTrailers
	WriterStateDone
)

type Writer struct {
	w           io.Writer
	writerState WriterState
}

func NewWriter(w io.Writer) Writer {
	return Writer{w: w, writerState: WriterStateInitialized}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	w.writerState = WriterStateWritingHeader
	return WriteStatusLine(w.w, statusCode)
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	w.writerState = WriterStateWritingBody
	return WriteHeaders(w.w, headers)
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	w.writerState = writerStateTrailers
	n, err := bytes.NewBuffer(p).WriteTo(w.w)
	return int(n), err
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	l := len(p)
	m, err := fmt.Fprintf(w.w, "%x\r\n", l)
	if err != nil {
		return m, err
	}
	n, err := bytes.NewBuffer(p).WriteTo(w.w)
	if err != nil {
		return m + int(n), err
	}
	r, err := fmt.Fprint(w.w, "\r\n")
	return m + int(n) + r, err
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.writerState != WriterStateWritingBody {
		return 0, fmt.Errorf("cannot write body in state %d", w.writerState)
	}
	defer func() { w.writerState = writerStateTrailers }()
	return fmt.Fprintf(w.w, "0\r\n")
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	if w.writerState != writerStateTrailers {
		return fmt.Errorf("cannot write trailers in state %d", w.writerState)
	}
	defer func() { w.writerState = WriterStateDone }()
	return w.WriteHeaders(h)
}
