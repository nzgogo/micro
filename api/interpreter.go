// NatsProxy serves as a proxy between gnats and http.

package gogoapi

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/nzgogo/micro/codec"
)

// HTTPReqToIntrlSReq creates the Request struct from regular *http.Request by serialization of main parts of it.
func HTTPReqToIntrlSReq(req *http.Request, rplSub, ctxid string) (*codec.Message, error) {
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
		Type:      "request",
		ContextID: ctxid,
		Header:    req.Header,
		Body:      string(buf.Bytes()),

		Method: req.Method,
		Path:   req.URL.Path,
		Host:   req.Host,

		ReplyTo: rplSub,
		Query:   req.URL.Query(),
		//Post
		//Scheme
	}
	return request, nil
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
