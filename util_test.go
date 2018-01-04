package gogo

import (
	"fmt"
	"testing"
)

func TestUrlReplace(t *testing.T) {
	path := "/home/:event/:session/:token"
	res := SubscribeURLToNats("POST", path)
	if res != "POST:.home.*.*.*" {
		fmt.Println(res)
		t.FailNow()
	}
}

