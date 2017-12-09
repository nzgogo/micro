package natsproxy

import (
	"bytes"
	"errors"
	"net/http"
	"micro/transport"
)

// Request wraps the HTTP request
// to be processed via pub/sub system.
//type Request struct {
//	URL        string
//	Method     string
//	Header     http.Header
//	Form       url.Values
//	RemoteAddr string
//	Body       []byte
//}

type Request *transport.Message

// NewRequestFromHTTP creates the Request struct from regular *http.Request by serialization of main parts of it.
func NewRequestFromHTTP(req *http.Request) (Request, error) {
	if req == nil {
		return nil, errors.New("natsproxy: Request cannot be nil")
	}
	var buf bytes.Buffer
	if req.Body != nil {
		//if err := req.ParseForm(); err != nil {
		//	return nil, err
		//}
		if _, err := buf.ReadFrom(req.Body); err != nil {
			return nil, err
		}
		if err := req.Body.Close(); err != nil {
			return nil, err
		}
	}

	request := Request{
		//URL:        req.URL.String(),
		//Method:     req.Method,
		Header:     make(map[string]string, 0),
		//Form:       req.Form,
		//RemoteAddr: req.RemoteAddr,
		Body:       buf.Bytes(),
	}
	return request, nil
}
