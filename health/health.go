package health

import (
	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type HealthResponse struct {
	ServiceName string
	ServiceId   string
	ServiceLoad int
}

func encode(hr *HealthResponse) []byte {
	return json.Marshal(hr)
}

func decode(res []byte) HealthResponse {
	var resStruct HealthResponse
	json.Unmarshal(res, resStruct)

	return resStruct
}
