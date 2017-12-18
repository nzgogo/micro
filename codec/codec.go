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
type Request struct {
	Method    string           `json:"method"`
	Path      string           `json:"path"`
	Authority string           `json:"authority"`
	Scheme    string           `json:"scheme"`
	Header    map[string]*Pair `json:"header"`
	Get       map[string]*Pair `json:"get"`
	Post      map[string]*Pair `json:"post"`
	Body      string           `json:"body"`
}

// Response struct represents a response message
type Response struct {
	StatusCode int              `json:"statusCode"`
	Header     map[string]*Pair `json:"header"`
	Body       string           `json:"body"`
}

// Codec is a interface
type Codec interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

type codec struct{}

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
