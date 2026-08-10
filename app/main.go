package main

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/http-server-starter-go/app/internal/http"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
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

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	req := http.NewRequest(conn)
	if req == nil {
		fmt.Println("failed to build the rquest: ")
		os.Exit(1)
	}

	response := "HTTP/1.1 200 OK\r\n\r\n"
	if req.Path != "/" {
		response = "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Println("an error corrued: ", err.Error())
		os.Exit(1)
	}
}
