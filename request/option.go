package request

import (
	"fmt"
	"httpclient/codec"
	"net/http"
	"net/url"
)

type Option func(*Request)

func New(method, u string, opts ...Option) *Request {
	r := &Request{
		Method: method,
		URL:    u,
		Header: make(http.Header),
		Query:  make(url.Values),
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func WithHeader(k, v string) Option {
	return func(r *Request) {
		r.Header.Add(k, v)
	}
}

func WithCodec(c codec.Codec) Option {
	return func(r *Request) {
		r.Codec = c
	}
}

func WithBody(body any) Option {
	return func(r *Request) {
		r.Body = body
	}
}

func WithURLParams(params map[string]string) Option {
	return func(r *Request) {
		for k, v := range params {
			r.Query.Add(k, v)
		}
	}
}

func WithURLParam(key string, value any) Option {
	return func(r *Request) {
		r.Query.Add(key, fmt.Sprint(value))
	}
}

func WithURLValues(values url.Values) Option {
	return func(r *Request) {
		for k, vs := range values {
			for _, v := range vs {
				r.Query.Add(k, v)
			}
		}
	}
}
