package main

import (
	"bytes"
	"io"
)

func main() {
	// listner, err := net.ResolveUDPAddr("localhost", ":42069")
	// if err != nil {
	// 	log.Fatal("error", "error", err)
	// }

	// net.DialUDP("")
	// for {
	// 	conn, err := listner.Accept()
	// 	if err != nil {
	// 		log.Fatal("error", "error", err)
	// 	}
	// 	lines := getLinesChannel(conn)
	// 	for line := range lines {
	// 		fmt.Println("read:", line)
	// 	}
	// }
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer f.Close()
		str := ""
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				break
			}
			data = data[:n]
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				str += string(data[:i])
				lines <- str
				str = string(data[i+1:])
			} else {
				str += string(data)
			}
		}
		if len(str) != 0 {
			lines <- str
		}
		close(lines)
	}()
	return lines
}
