package http

import (
	"fmt"
)

type Response struct {
	Version    string
	Code       string
	CodeStatus string
	Header     Header
	Body       any
}

func Response404(version string) Response {
	return Response{
		Version:    version,
		Code:       "404",
		CodeStatus: "Not Found",
	}
}

func (r Response) String() string {
	str := r.Version + " " + r.Code + " " + r.CodeStatus + CRLF
	for key, header := range r.Header {
		str += fmt.Sprintf("%s: %v%s", key, header, CRLF)
	}
	str += CRLF
	if r.Body != nil {
		str += fmt.Sprintf("%v", r.Body)
	}
	return str
}
