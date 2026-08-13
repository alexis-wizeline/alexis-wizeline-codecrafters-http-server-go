package http

import (
	"fmt"
	"net"
	"os"
)

const (
	CRLF = "\r\n"
)

type HandleFunc func(r *Request) Response
type Routes map[string]HandleFunc

type Server struct {
	routes    Routes
	folder    string
	directory string
}

func NewServer(directory string) *Server {
	return &Server{
		routes:    make(Routes),
		directory: directory,
	}
}

func (s *Server) Directory() string {
	return s.directory
}

func (s *Server) Listen(network, address string) (net.Listener, error) {
	return net.Listen(network, address)
}

func (s *Server) HandleConn(conn net.Conn) {
	defer conn.Close()
	req := NewRequest(conn)
	if req == nil {
		fmt.Println("failed to build the rquest: ")
		os.Exit(1)
	}

	reqParams, handler := s.GetRouteHandler(req.Path)
	req.AddParams(reqParams)
	response := Response404(req.Version)
	if handler != nil {
		response = handler(req)
	}
	_, err := conn.Write([]byte(response.String()))
	if err != nil {
		fmt.Println("an error corrued: ", err.Error())
		os.Exit(1)
	}
}

func (s *Server) Register(pattern string, handler HandleFunc) {
	s.routes[pattern] = handler
}

func (s *Server) GetRouteHandler(path string) (Params, HandleFunc) {
	for patter, handler := range s.routes {
		match, params := matchPath(patter, path)
		if match {
			return params, handler
		}
	}

	return nil, nil
}

func matchPath(pattern, path string) (bool, Params) {
	patternSegs := segmentBy(pattern, '/')
	pathSegs := segmentBy(path, '/')
	if len(patternSegs) != len(pathSegs) {
		return false, nil
	}
	params := make(Params)
	for i, seg := range patternSegs {
		if isParam(seg) {
			key := seg[1 : len(seg)-1]
			params[key] = pathSegs[i]
		} else if seg != pathSegs[i] {
			return false, nil
		}

	}
	return true, params
}

func segmentBy(s string, delimeter byte) []string {
	segments := make([]string, 0)
	segment := ""
	for _, c := range s {
		if c == rune(delimeter) {
			if len(segment) > 0 {
				segments = append(segments, segment)
				segment = ""
			}
			continue
		}
		segment += string(c)
	}

	if len(segment) > 0 {
		segments = append(segments, segment)
	}

	return segments
}

func isParam(segment string) bool {
	if len(segment) == 0 {
		return false
	}
	last := len(segment) - 1
	return segment[0] == '{' && segment[last] == '}'
}
