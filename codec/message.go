package codec

import (
	"net/http"

	validator "github.com/asaskevich/govalidator"
)

type Message struct {
	//HTTP request mandatory fields
	Method string `json:"method,omitempty"`
	Path   string `json:"path,omitempty"`
	Host   string `json:"host,omitempty"`
	Scheme string `json:"scheme,omitempty"`

	//Internal request fields
	ReplyTo   string `json:"replyTo,omitempty"`
	Node      string `json:"node,omitempty"`
	ContextID string `json:"contextID,omitempty"`

	//Internal response fields
	StatusCode int `json:"statusCode,omitempty"`

	//Common fields
	Type   string                 `json:"type,omitempty"`
	Header http.Header            `json:header,omitempty`
	Body   map[string]interface{} `json:body,omitempty`
}

func (msg *Message) Set(key string, value interface{}) {
	msg.Body[key] = value
}

func (msg *Message) Del(key string) {
	delete(msg.Body, key)
}

func (msg *Message) Get(key string) (value interface{}, ok bool) {
	value, ok = msg.Body[key]
	return
}

func (msg *Message) GetString(key string) (value string, ok bool) {
	v, ok := msg.Body[key]
	if !ok {
		return
	}

	value = validator.ToString(v)

	return
}

func (msg *Message) GetInt(key string) (value int64, ok bool) {
	v, ok := msg.Body[key]
	if !ok {
		return
	}

	stringValue := validator.ToString(v)
	value, err := validator.ToInt(stringValue)
	if err != nil {
		ok = false
		return
	}

	return
}

func (msg *Message) GetFloat(key string) (value float64, ok bool) {
	v, ok := msg.Body[key]
	if !ok {
		return
	}

	stringValue := validator.ToString(v)
	value, err := validator.ToFloat(stringValue)
	if err != nil {
		ok = false
		return
	}

	return
}

func (msg *Message) GetBool(key string) (value bool, ok bool) {
	v, ok := msg.Body[key]
	if !ok {
		return
	}

	stringValue := validator.ToString(v)
	value, err := validator.ToBoolean(stringValue)
	if err != nil {
		ok = false
		return
	}

	return
}

func (msg *Message) ParseHTTPRequest(r *http.Request) *Message {
	return msg
}

func (msg *Message) WriteHTTPResponse(rw http.ResponseWriter) {}
