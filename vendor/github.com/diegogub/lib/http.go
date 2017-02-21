package lib

import (
	"net/http"
)

func NewHttpHeader() *http.Header {
	var h http.Header
	h = make(map[string][]string)
	return &h
}
