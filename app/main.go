package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/http-server-starter-go/app/internal/http"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	s := http.NewServer()
	s.Register("/", func(r *http.Request) http.Response {
		return http.Response{
			Version:    r.Version,
			Code:       "200",
			CodeStatus: "OK",
		}
	})
	s.Register("/echo/{str}", func(r *http.Request) http.Response {
		headers := make(http.Header)
		body, ok := r.Params["str"].(string)
		if ok {
			headers["Content-Type"] = "text/plain"
			headers["Content-Length"] = len(body)
		}
		return http.Response{
			Version:    r.Version,
			Code:       "200",
			CodeStatus: "OK",
			Header:     headers,
			Body:       body,
		}
	})

	l, err := s.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go s.HandleConn(conn)
	}
}
