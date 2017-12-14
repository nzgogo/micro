package codec
<<<<<<< HEAD

// Codec interface
=======
// Codec is used for encoding where the transport doesn't natively support
// headers in the json type. In this case the entire message is
// encoded as the payload

import "encoding/json"

>>>>>>> fb9de2dd328134ab8842aa1f013bed95261e9fde
type Codec interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
	String() string
}
<<<<<<< HEAD
=======
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
>>>>>>> fb9de2dd328134ab8842aa1f013bed95261e9fde
