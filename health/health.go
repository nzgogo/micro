package health

import (
	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// ResponseMsg struct
type ResponseMsg struct {
	ServiceName string
	ServiceID   string
	ServiceLoad int
}

// Encode function
func Encode(hr *ResponseMsg) ([]byte, error) {
	return json.Marshal(hr)
}

// Decode function
func Decode(res []byte) *ResponseMsg {
	var resStruct ResponseMsg
	json.Unmarshal(res, resStruct)

	return &resStruct
}
