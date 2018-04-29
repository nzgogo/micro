package codec

import (
	"net/http"

	"github.com/json-iterator/go"
	"github.com/nzgogo/micro/constant"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func Unmarshal(d []byte, v interface{}) error {
	return json.Unmarshal(d, v)
}

func NewMessage(t string) *Message {
	return &Message{
		Type:   t,
		Header: http.Header{},
		Body:   make(map[string]interface{}),
	}
}

//NewResponse creates Response Message object.
func NewResponse(contextID string, statusCode int) *Message {
	return &Message{
		Type:       constant.RESPONSE,
		StatusCode: statusCode,
		ContextID:  contextID,
		Header:     http.Header{},
		Body:       make(map[string]interface{}),
	}
}

func NewJsonResponse(contextID string, statusCode int) *Message {
	h := http.Header{}
	h.Add("Content-Type", "application/json")

	return &Message{
		Type:       constant.RESPONSE,
		StatusCode: statusCode,
		ContextID:  contextID,
		Header:     h,
		Body:       make(map[string]interface{}),
	}
}
