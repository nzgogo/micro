// NatsProxy serves as a proxy between gnats and http.

package gogoapi

import (
	"bytes"
	"net/http"
	"errors"
	//"github.com/nzgogo/micro/codec"
	"micro/codec"
)

// Request Response server as structure to transport http response throu NATS message queue
type Request *codec.Request
type Response *codec.Response

// NewRequestFromHTTP creates the Request struct from regular *http.Request by serialization of main parts of it.
func HTTPReqToNatsSReq(req *http.Request) (Request, error) {
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

	request := &codec.Request{
		Method:     req.Method,
		Path:		req.RequestURI,
		Authority:	req.Host,
		Body:       string(buf.Bytes()),
	}
	return request, nil
}

// NewResponse creates blank initialized Response object.
func NewResponse() Response {
	return &codec.Response{
		200,
		make(map[string][]string, 0),
		"",
	}
}

func WriteResponse(rw http.ResponseWriter, response Response) {
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


