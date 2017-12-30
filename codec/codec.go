package codec

import (
	"github.com/json-iterator/go"
)

// Codec is a interface
type Codec interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

type codec struct{}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Pair struct represents a key-value pair
type Pair struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

// Request struct represents a request message
type Request struct {
	Method string              `json:"method,omitempty"`
	Path   string              `json:"path,omitempty"`
	Host   string              `json:"host,omitempty"`
	Scheme string              `json:"scheme"`
	Node   string              `json:"node,omitempty"`
	Header map[string][]string `json:"header"`
	Get    map[string]*Pair    `json:"get,omitempty"`
	Post   map[string]*Pair    `json:"post,omitempty"`
	Body   string              `json:"body"`
}

// Response struct represents a response message
type Response struct {
	StatusCode int                 `json:"statusCode"`
	Header     map[string][]string `json:"header"`
	Body       string              `json:"body"`
}

func (j codec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j codec) Unmarshal(d []byte, v interface{}) error {
	return json.Unmarshal(d, v)
}

// NewCodec returns a new json codec
func NewCodec() Codec {
	return codec{}
}
