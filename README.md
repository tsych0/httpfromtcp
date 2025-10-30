# HTTP from TCP

A low-level HTTP/1.1 server implementation built from scratch using raw TCP connections in Go. This project is part of Boot.dev's "HTTP from TCP" course, demonstrating how HTTP works at the protocol level without using high-level HTTP libraries.

## 🎯 What This Project Does

This project implements a fully functional HTTP/1.1 server by parsing TCP byte streams directly. It handles:
- HTTP request parsing (request line, headers, body)
- HTTP response generation (status lines, headers, body)
- Chunked transfer encoding
- HTTP trailers
- Multiple endpoints with different response types (HTML, video, proxy)

## 🚀 Features

- **Custom HTTP Parser**: Parses HTTP requests byte-by-byte from TCP connections
- **Request Handler System**: Flexible handler pattern for routing requests
- **Multiple Content Types**: Serves HTML, plain text, and video files
- **Chunked Encoding**: Implements chunked transfer encoding with trailers
- **HTTP Proxy**: Proxies requests to httpbin.org with SHA256 content verification
- **Graceful Shutdown**: Handles SIGINT/SIGTERM signals properly

## 📁 Project Structure

```
.
├── cmd/
│   ├── httpserver/      # Main HTTP server implementation
│   ├── tcplistener/     # Simple TCP listener for debugging
│   └── udpsender/       # UDP utilities (experimental)
├── internal/
│   ├── headers/         # HTTP headers parser
│   ├── request/         # HTTP request parser
│   ├── response/        # HTTP response writer
│   └── server/          # TCP server and handler logic
├── assets/              # Static files (video)
└── messages.txt         # Sample text file
```

## 🛠️ Running the Server

```bash
# Run the HTTP server
go run cmd/httpserver/*.go

# The server starts on port 42069
# Visit http://localhost:42069 in your browser
```

## 📡 Available Endpoints

- `GET /` - Success page (200 OK)
- `GET /yourproblem` - Bad request page (400)
- `GET /myproblem` - Server error page (500)
- `GET /video` - Serves a video file (MP4)
- `GET /httpbin/*` - Proxies to httpbin.org with chunked encoding and SHA256 trailers

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Test specific packages
go test ./internal/request
go test ./internal/headers
```

### Test with curl

```bash
# Basic GET request
curl http://localhost:42069/

# POST request with body
curl -X POST http://localhost:42069/ -d "hello world"

# Proxy request with verbose output
curl -v http://localhost:42069/httpbin/get
```

## 🧠 Key Concepts Learned

1. **TCP Sockets**: Direct communication using `net.Listen` and `net.Conn`
2. **HTTP Protocol**: Understanding HTTP/1.1 request/response format
3. **Byte Stream Parsing**: Building a parser that handles partial reads
4. **State Machines**: Implementing stateful parsing logic
5. **Chunked Transfer Encoding**: Streaming responses without knowing content length
6. **HTTP Trailers**: Sending headers after the body (e.g., checksums)
7. **Concurrent Connections**: Handling multiple clients simultaneously

## 📝 Implementation Highlights

### Custom Request Parser
- Handles variable-sized reads from TCP stream
- Parses request line, headers, and body sequentially
- Validates HTTP/1.1 format and header field names

### Response Writer
- State machine ensuring correct write order (status → headers → body → trailers)
- Supports both regular and chunked response bodies
- Automatic content-length calculation

### Chunked Encoding with Trailers
- Streams data in chunks without knowing total size upfront
- Computes SHA256 hash of proxied content
- Sends hash as trailer after body completion

## 🔧 Technologies Used

- **Go 1.24.5**
- **net package**: Raw TCP socket programming
- **crypto/sha256**: Content hashing
- **testing**: Unit tests with table-driven approach

## 🎓 Course Credit

This project was completed as part of [Boot.dev's HTTP from TCP course](https://boot.dev), which teaches low-level networking and protocol implementation.

## 👤 Author

**Ayush Biswas** ([@tsych0](https://github.com/tsych0))

---

*Note: This server is for educational purposes and should not be used in production. It lacks many features required for a production HTTP server (TLS, HTTP/2, connection pooling, etc.).*
