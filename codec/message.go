package codec

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"

	validator "github.com/asaskevich/govalidator"
	"github.com/nzgogo/micro/constant"
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
	Type    string                 `json:"type,omitempty"`
	Header  http.Header            `json:"header,omitempty"`
	RawBody []byte                 `json:"rawBody,omitempty"`
	Body    map[string]interface{} `json:"body,omitempty"`
}

func NewMessage(t string) *Message {
	return &Message{
		Type: t,
		Body: make(map[string]interface{}),
	}
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

func (msg *Message) GetBytes(key string) (value []byte, ok bool) {
	v, o := msg.Body[key]
	if !o {
		return
	}
	value, ok = v.([]byte)

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

func (msg *Message) ParseHTTPRequest(r *http.Request, replyTo string, contextID string) (*Message, error) {
	msg.Body = make(map[string]interface{})

	r.ParseForm()
	for k, v := range r.Form {
		if len(v) == 1 {
			msg.Body[k] = v[0]
			continue
		}
		msg.Body[k] = v
	}

	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			return nil, err
		}
		msg.RawBody = b
		var j map[string]interface{}
		Unmarshal(b, &j)
		for k, v := range j {
			msg.Body[k] = v
		}
	} else if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		r.ParseMultipartForm(0)
		for k, v := range r.MultipartForm.Value {
			if len(v) == 1 {
				msg.Body[k] = v[0]
				continue
			}
			msg.Body[k] = v
		}
		file, fileHeader, err := r.FormFile("file")
		if err == nil {
			fileRaw := make([]byte, fileHeader.Size)
			file.Read(fileRaw)
			msg.Body["file"] = fileRaw
		}
	}

	msg.Type = constant.REQUEST
	msg.ContextID = contextID
	msg.ReplyTo = replyTo
	msg.Method = r.Method
	msg.Host = r.Host
	msg.Path = r.URL.Path
	msg.Header = r.Header
	return msg, nil
}

func (msg *Message) WriteHTTPResponse(rw http.ResponseWriter, response *Message) {
	for k, values := range msg.Header {
		for _, v := range values {
			rw.Header().Add(k, v)
		}
	}
	rw.WriteHeader(msg.StatusCode)
	if b, err := Marshal(msg.Body); err == nil {
		bytes.NewBuffer(b).WriteTo(rw)
	}
}
