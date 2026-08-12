package http

import (
	"bufio"
	"io"
	"maps"
	"net"
	"strings"
	"time"
)

type Params map[string]any
type Request struct {
	Method  string
	Path    string
	Version string
	Header  Header
	Params  Params
}

func (req *Request) AddParams(new Params) {
	if new == nil {
		return
	}
	if req.Params == nil {
		req.Params = make(Params)
	}

	maps.Copy(req.Params, new)
}

func NewRequest(conn net.Conn) *Request {
	stringReq, err := readConn(conn)
	if err != nil {
		return nil
	}

	requestLine, headersSlice := cleanStr(cleanStr(stringReq[0], '\r'), '\n'), stringReq[1:]
	requestParts := segmentBy(requestLine, ' ')
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

		headers[key.String()] = cleanStr(cleanStr(cleanStr(header[index+1:], '\r'), '\n'), ' ')
	}

	path, params := parsePath(requestParts[1])
	req := &Request{
		Method:  requestParts[0],
		Path:    path,
		Version: requestParts[2],
		Header:  headers,
		Params:  params,
	}
	return req
}

func readConn(conn net.Conn) ([]string, error) {
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	var lines []string
	readder := bufio.NewReader(conn)
	for {
		line, err := readder.ReadString('\n')
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

func parsePath(path string) (string, Params) {
	params := make(Params)
	qParamsStartIndex := -1
	for index, char := range path {
		if char == '?' {
			qParamsStartIndex = index
			break
		}
	}
	if qParamsStartIndex == -1 {
		return path, params
	}

	queryParams := path[qParamsStartIndex+1:]
	key, value := "", ""
	isKey := true
	for _, char := range queryParams {
		if isKey {
			if char == '=' {
				isKey = false
				continue
			}
			key += string(char)
		} else {
			if char == '&' {
				params[key] = value
				key, value = "", ""
				isKey = true
				continue
			}
			value += string(char)
		}
	}
	params[key] = value

	pathWithoutQuery := path[0:qParamsStartIndex]
	return pathWithoutQuery, params
}

func cleanStr(str string, rm byte) string {
	res := strings.Builder{}
	for _, ch := range str {
		if ch == rune(rm) {
			continue
		}
		res.WriteByte(byte(ch))
	}

	return res.String()
}
