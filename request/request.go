package request

import (
	"httpclient/codec"
	"net/http"
	"net/url"
)

type Request struct {
	Method string
	URL    string
	Header http.Header
	Body   any
	Codec  codec.Codec
	Query  url.Values
}
