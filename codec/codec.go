package codec
// Codec is used for encoding where the transport doesn't natively support
// headers in the json type. In this case the entire message is
// encoded as the payload

import "encoding/json"

type Codec interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
	String() string
}
type jsonCodec struct{}

func (j jsonCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j jsonCodec) Unmarshal(d []byte, v interface{}) error {
	return json.Unmarshal(d, v)
}

func (j jsonCodec) String() string {
	return "json"
}

func NewCodec() Codec {
	return jsonCodec{}
}
