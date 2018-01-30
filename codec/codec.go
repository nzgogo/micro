package codec

import (
	"net/http"
	"net/url"

	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Pair struct represents a key-value pair
type Pair struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

type Message struct {
	//pre-process
	Method string `json:"method,omitempty"`
	Path   string `json:"path,omitempty"`
	Host   string `json:"host,omitempty"`

	//request fields
	ReplyTo string           `json:"replyTo,omitempty"`
	Node    string           `json:"node,omitempty"`
	Query   url.Values       `json:"get,omitempty"`
	Post    map[string]*Pair `json:"post,omitempty"`
	Scheme  string           `json:"scheme"`

	//response fields
	StatusCode int `json:"statusCode"`

	//common fields
	Type      string      `json:"type"`
	ContextID string      `json:"contextId"`
	Header    http.Header `json:"header"`
	Body      []byte      `json:"body"`
}

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func Unmarshal(d []byte, v interface{}) error {
	return json.Unmarshal(d, v)
}

//NewResponse creates Response Message object.
func NewResponse(statusCode int, contextID string, body *string, header http.Header) *Message {
	var b = make([]byte, 0)
	if body != nil {
		b = []byte(*body)
	} else {
		b = nil
	}

	return &Message{
		Type:       "response",
		StatusCode: statusCode,
		Header:     header,
		ContextID:  contextID,
		Body:       b,
	}
}

func (msg *Message) Get(key string) string {
	jsonStrings := make(map[string]string)
	if err := Unmarshal(msg.Body, &jsonStrings); err == nil {
		return jsonStrings[key]
	}
	return ""
}

func (msg *Message) GetAll(key string) map[string]string {
	jsonStrings := make(map[string]string, 0)
	err := Unmarshal(msg.Body, &jsonStrings)
	if err == nil {
		return jsonStrings
	}
	return nil
}
