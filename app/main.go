package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/codecrafters-io/http-server-starter-go/app/internal/files"
	"github.com/codecrafters-io/http-server-starter-go/app/internal/http"
)

var directory = flag.String("directory", "", "the directory of the server")

func main() {
	flag.Parse()
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	s := http.NewServer(*directory)
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
	s.Register("/user-agent", func(r *http.Request) http.Response {
		userAgent, ok := r.Header["User-Agent"].(string)
		headers := make(http.Header)
		if ok {
			headers["Content-Type"] = "text/plain"
			headers["Content-Length"] = len(userAgent)
		}
		return http.Response{
			Version:    r.Version,
			Code:       "200",
			CodeStatus: "OK",
			Header:     headers,
			Body:       userAgent,
		}
	})

	s.Register("/files/{filename}", func(r *http.Request) http.Response {
		fileName, ok := r.Params["filename"].(string)
		if !ok || len(fileName) == 0 {
			return http.Response{
				Version:    r.Version,
				Code:       "400",
				CodeStatus: "Bad Request",
			}
		}
		content, err := files.ReadFile(s.Directory(), fileName)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return http.Response404(r.Version)
			}
		}

		headers := make(http.Header)
		headers["Content-Type"] = "application/octet-stream"
		headers["Content-Length"] = len(content)

		return http.Response{
			Version:    r.Version,
			Code:       "200",
			CodeStatus: "OK",
			Header:     headers,
			Body:       string(content),
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
