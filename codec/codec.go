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

func (msg *Message) Get(key string) string {
	jsonStrings := make(map[string]string)
	if err := Unmarshal(msg.Body, &jsonStrings); err == nil {
		return jsonStrings[key]
	}
	return ""
}

func (msg *Message) GetAll() map[string]string {
	jsonStrings := make(map[string]string, 0)
	err := Unmarshal(msg.Body, &jsonStrings)
	if err == nil {
		return jsonStrings
	}
	return nil
}

//NewResponse creates Response Message object.
func NewResponse(contextID string, statusCode int, body []byte, header http.Header) *Message {
	return &Message{
		Type:       "response",
		StatusCode: statusCode,
		Header:     header,
		ContextID:  contextID,
		Body:       body,
	}
}

func NewJsonResponse(contextID string, statusCode int, body interface{}) *Message {
	var b []byte
	if v, ok := body.([]byte); ok {
		b = v
	} else {
		if v, err := Marshal(body); err != nil {
			b = nil
		} else {
			b = v
		}
	}

	h := http.Header{}
	h.Add("Content-Type", "application/json")
	return &Message{
		Type:       "response",
		StatusCode: statusCode,
		Header:     h,
		ContextID:  contextID,
		Body:       b,
	}
}
