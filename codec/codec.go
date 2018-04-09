package codec

import (
	"net/http"

	"github.com/json-iterator/go"
	"github.com/nzgogo/micro/constant"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Pair struct represents a key-value pair
// type Pair struct {
// 	Key    string   `json:"key"`
// 	Values []string `json:"values"`
// }

// type Message struct {
// 	//pre-process
// 	Method string `json:"method,omitempty"`
// 	Path   string `json:"path,omitempty"`
// 	Host   string `json:"host,omitempty"`
//
// 	//request fields
// 	ReplyTo string           `json:"replyTo,omitempty"`
// 	Node    string           `json:"node,omitempty"`
// 	Query   url.Values       `json:"get,omitempty"`
// 	Post    map[string]*Pair `json:"post,omitempty"`
// 	Scheme  string           `json:"scheme"`
//
// 	//response fields
// 	StatusCode int `json:"statusCode"`
//
// 	//common fields
// 	Type      string      `json:"type"`
// 	ContextID string      `json:"contextID"`
// 	Header    http.Header `json:"header"`
// 	Body      []byte      `json:"body"`
// }

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func Unmarshal(d []byte, v interface{}) error {
	return json.Unmarshal(d, v)
}

// func (msg *Message) Get(key string) interface{} {
// 	jsonStrings := make(map[string]interface{})
// 	if err := Unmarshal(msg.Body, &jsonStrings); err == nil {
// 		return jsonStrings[key]
// 	}
// 	return nil
// }
//
// func (msg *Message) GetAll() map[string]interface{} {
// 	jsonStrings := make(map[string]interface{}, 0)
// 	err := Unmarshal(msg.Body, &jsonStrings)
// 	if err == nil {
// 		return jsonStrings
// 	}
// 	return nil
// }
//
// func (msg *Message) Set(key string, value interface{}) {
// 	body := make(map[string]interface{}, 0)
// 	Unmarshal(msg.Body, &body)
// 	body[key] = value
// 	newMsg, err := Marshal(body)
// 	if err == nil {
// 		msg.Body = newMsg
// 	}
// }

//NewResponse creates Response Message object.
func NewResponse(contextID string, statusCode int, body map[string]interface{}, header http.Header) *Message {
	return &Message{
		Type:       constant.RESPONSE,
		StatusCode: statusCode,
		Header:     header,
		ContextID:  contextID,
		Body:       body,
	}
}

func NewJsonResponse(contextID string, statusCode int, body map[string]interface{}) *Message {
	h := http.Header{}
	h.Add("Content-Type", "application/json")

	return &Message{
		Type:       constant.RESPONSE,
		StatusCode: statusCode,
		Header:     h,
		ContextID:  contextID,
		Body:       body,
	}
}
