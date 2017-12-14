package json

import (
	"github.com/json-iterator/go"
	"github.com/nzgogo/micro/codec"
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

func NewCodec() codec.Codec {
	return jsonCodec{}
}
