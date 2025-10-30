package main

import (
	"os"

	"github.com/tsych0/httpfromtcp/internal/request"
	"github.com/tsych0/httpfromtcp/internal/response"
)

const badRequestHTML = `
<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
`

func yourProblem(w *response.Writer, r *request.Request) {
	w.WriteStatusLine(response.StatusBadRequest)
	headers := response.GetDefaultHeaders(len(badRequestHTML))
	headers.Set("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody([]byte(badRequestHTML))
}

const serverErrorHTML = `
<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
`

func myProblem(w *response.Writer, r *request.Request) {
	w.WriteStatusLine(response.StatusInternalServerError)
	headers := response.GetDefaultHeaders(len(serverErrorHTML))
	headers.Set("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody([]byte(serverErrorHTML))
}

const bangerHTML = `
<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
`

func banger(w *response.Writer, r *request.Request) {
	w.WriteStatusLine(response.StatusOk)
	headers := response.GetDefaultHeaders(len(bangerHTML))
	headers.Set("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody([]byte(bangerHTML))
}

func giveTheVideo(w *response.Writer, r *request.Request) {
	w.WriteStatusLine(response.StatusOk)
	file, err := os.ReadFile("assets/vim.mp4")
	if err != nil {
		myProblem(w, r)
		return
	}
	headers := response.GetDefaultHeaders(len(file))
	headers.Set("Content-Type", "video/mp4")
	w.WriteHeaders(headers)
	w.WriteBody(file)
}
