package codec

import (
	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Pair struct represents a key-value pair
type Pair struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

// Request struct represents a request message
// type Request struct {
// 	Method string              `json:"method,omitempty"`
// 	Path   string              `json:"path,omitempty"`
// 	Host   string              `json:"host,omitempty"`
// 	Scheme string              `json:"scheme"`
// 	Node   []byte              `json:"node,omitempty"`
// 	Header map[string][]string `json:"header"`
// 	Get    map[string]*Pair    `json:"get,omitempty"`
// 	Post   map[string]*Pair    `json:"post,omitempty"`
// 	Body   string              `json:"body"`
// }

// Response struct represents a response message
// type Response struct {
// 	StatusCode int                 `json:"statusCode"`
// 	Header     map[string][]string `json:"header"`
// 	Body       string              `json:"body"`
// }

type Message struct {
	//pre-process
	Method string `json:"method,omitempty"`
	Path   string `json:"path,omitempty"`
	Host   string `json:"host,omitempty"`

	//request fields
	ReplyTo string              `json:"replyTo,omitempty"`
	Node    string              `json:"node,omitempty"`
	Query   map[string][]string `json:"get,omitempty"`
	Post    map[string]*Pair    `json:"post,omitempty"`
	Scheme  string              `json:"scheme"`

	//response fields
	StatusCode int `json:"statusCode"`

	//common fields
	Type      string              `json:"type"`
	ContextID string              `json:"contextId"`
	Header    map[string][]string `json:"header"`
	Body      string              `json:"body"`
}

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func Unmarshal(d []byte, v interface{}) error {
	return json.Unmarshal(d, v)
}
