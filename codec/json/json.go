package json

import (
	"github.com/json-iterator/go"
	bCodec "github.com/micro/go-micro/broker/codec"
	tCodec "github.com/micro/go-micro/transport/codec"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

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

// NewBrokerCodec returns a go-micro broker Codec interface compatible new json codec
func NewBrokerCodec() bCodec.Codec {
	return jsonCodec{}
}

// NewTransportCodec returns a go-micro transport Codec interface compatible new json codec
func NewTransportCodec() tCodec.Codec {
	return jsonCodec{}
}
