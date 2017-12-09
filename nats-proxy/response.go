package natsproxy

import (
	"micro/transport"
)

// Response server as structure
// to transport http response
// throu NATS message queue
//type Response struct {
//	Header     http.Header
//	StatusCode int
//	Body       []byte
//}
type Response *transport.Message
// NewResponse creates blank
// initialized Response object.
func NewResponse() Response {
	return &transport.Message{
		make(map[string]string, 0),
		//200,
		make([]byte, 0),
	}
}

