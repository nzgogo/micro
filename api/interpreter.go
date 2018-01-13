// NatsProxy serves as a proxy between gnats and http.

package gogoapi

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/nzgogo/micro/codec"
)

// NewRequestFromHTTP creates the Request struct from regular *http.Request by serialization of main parts of it.
func HTTPReqToNatsSReq(req *http.Request) (*codec.Message, error) {
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

	//TODO May need extract more data from http reqeust
	request := &codec.Message{
		Method: req.Method,
		Path:   req.RequestURI,
		Host:   req.Host,
		Body:   string(buf.Bytes()),
	}
	return request, nil
}

// NewResponse creates blank initialized Response object.
func NewResponse() *codec.Message {
	return &codec.Message{
		StatusCode: 200,
		Header:     make(map[string][]string, 0),
		Body:       "",
	}
}

func WriteResponse(rw http.ResponseWriter, response *codec.Message) {
	// Copy headers
	// from NATS response.
	copyHeader(response.Header, rw.Header())
	statusCode := response.StatusCode
	// Write the response code
	rw.WriteHeader(statusCode)

	// Write the bytes of response to a response writer.
	// TODO benchmark
	bytes.NewBuffer([]byte(response.Body)).WriteTo(rw)
}

func copyHeader(src, dst http.Header) {
	for k, v := range src {
		for _, val := range v {
			dst.Add(k, val)
		}
	}
}
