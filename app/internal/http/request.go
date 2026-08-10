package http

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type Request struct {
	Method string
	Path   string
	Versio string
	Header Header
}

func NewRequest(conn net.Conn) *Request {
	stringReq, err := readConn(conn)
	if err != nil {
		return nil
	}

	requestLine, headersSlice := stringReq[0], stringReq[1:]
	requestParts := strings.Split(requestLine, " ")
	if len(requestParts) < 3 {
		return nil
	}
	headers := make(Header)
	for _, header := range headersSlice {
		key := strings.Builder{}
		index := 0
		for {
			if header[index] == ':' {
				break
			}
			key.WriteByte(header[index])
			index++
		}

		headers[key.String()] = header[index+1:]
	}

	req := &Request{
		Method: requestParts[0],
		Path:   requestParts[1],
		Versio: requestParts[2],
		Header: headers,
	}

	return req
}

func readConn(conn net.Conn) ([]string, error) {
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	var lines []string
	readder := bufio.NewReader(conn)
	for {
		line, err := readder.ReadString('\n')
		fmt.Println(line)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if line == "\r\n" {
			break
		}
		lines = append(lines, line)
	}

	return lines, nil
}
