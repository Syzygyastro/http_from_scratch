package main

import (
	"fmt"
	"httpFromScratch/internal/request"
	"log"
	"net"
)

func hasNewline(b []byte) bool {
	for _, c := range b {
		if c == 10 {
			return true
		}
	}
	return false
}

// func getLinesChannel(f io.ReadCloser) <-chan string {
// 	output := make(chan string)
// 	go func() {
// 		defer f.Close()
// 		defer close(output)
// 		b := make([]byte, 8)
// 		line := ""
// 		for {
// 			n, err := f.Read(b)
// 			word := string(b[:n])
// 			if hasNewline(b[:n]) {
// 				parts := strings.Split(word, "\n")
// 				line += parts[0]
// 				output <- line
// 				if len(parts) == 2 {
// 					line = parts[1]
// 				} else {
// 					line = ""
// 				}

// 			} else {
// 				line += word
// 			}
// 			if err != nil {
// 				if line != "" {
// 					output <- line // Flush remaining data
// 				}
// 				if errors.Is(err, io.EOF) {
// 					break // Graceful exit
// 				}
// 				fmt.Printf("error: %s\n", err.Error()) // Log other errors
// 				return
// 			}
// 		}
// 	}()
// 	return output
// }

func main() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("=============================")
		fmt.Println("Server Listening on Port 42069")

		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println("Error in parsing input data", err.Error())
		}
		// Print the parsed request line
		fmt.Printf("Request line:\n")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)

		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}

}
