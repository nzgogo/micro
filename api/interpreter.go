// NatsProxy serves as a proxy between gnats and http.

package gogoapi

import (
	"bytes"
	"net/http"

	"github.com/nzgogo/micro/codec"
	"github.com/nzgogo/micro/constant"
)

// HTTPReqToIntrlSReq creates the Request struct from regular *http.Request by serialization of main parts of it.
func HTTPReqToIntrlSReq(req *http.Request, rplSub, ctxid string) (*codec.Message, error) {
	if req == nil {
		return nil, constant.ErrHttpEmptyRequest
	}
	var buf bytes.Buffer
	if req.Body != nil {
		if _, err := buf.ReadFrom(req.Body); err != nil {
			return nil, err
		}
		if err := req.Body.Close(); err != nil {
			return nil, err
		}
	}

	//TODO May need extract more data from http request
	request := &codec.Message{
		Type:      constant.REQUEST,
		ContextID: ctxid,
		Header:    req.Header,
		Body:      buf.Bytes(),

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
	// Copy headers from NATS response.
	if response == nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	copyHeader(response.Header, rw.Header())
	statusCode := response.StatusCode
	// Write the response code
	rw.WriteHeader(statusCode)

	// Write the bytes of response to a response writer.
	bytes.NewBuffer(response.Body).WriteTo(rw)
}

func copyHeader(src, dst http.Header) {
	if src == nil {
		return
	}
	for k, v := range src {
		for _, val := range v {
			dst.Add(k, val)
		}
	}
}
