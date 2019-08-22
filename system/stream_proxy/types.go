package stream_proxy

import (
	"net/http"
)

type StreamRequestModel struct {
	URI    string      `json:"uri"`
	Method string      `json:"method"`
	Body   []byte      `json:"body"`
	Header http.Header `json:"header"`
}

type StreamResponseModel struct {
	Code   int         `json:"code"`
	Body   []byte      `json:"body"`
	Header http.Header `json:"header"`
}
