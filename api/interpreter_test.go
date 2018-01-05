package gogoapi

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestHTTPReqToNatsSReq(t *testing.T) {
	url, _ := url.Parse("http://test.com/test")
	httpReq := &http.Request{
		Method: "GET",
		URL:    url,
		Body:   ioutil.NopCloser(bytes.NewReader([]byte{0xFF, 0xFC})),
	}
	req, err := HTTPReqToNatsSReq(httpReq)

	if err != nil {
		t.Error(err)
	}

	// if req.Host != "http://test.com" {
	// 	t.Error("Url not equals")
	// }
	if len(req.Body) != 2 {
		t.Error("Body length not equals")
	}
}
