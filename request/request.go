package request

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/imattdu/httpclient/codec"
)

type Request struct {
	Method  string
	URL     string
	Query   url.Values
	Header  http.Header
	Body    any
	RawBody []byte

	Codec  codec.Codec
	Stream bool // ⭐ 是否流式
}

var (
	ErrInvalidMethod = errors.New("request method is empty")
	ErrInvalidURL    = errors.New("request url is empty")
	ErrNoCodec       = errors.New("request body provided but codec is nil")
)

func (r *Request) Validate() error {
	if r.Method == "" {
		return ErrInvalidMethod
	}
	if r.URL == "" {
		return ErrInvalidURL
	}
	if r.Body != nil && r.Codec == nil {
		return ErrNoCodec
	}
	return nil
}

func (r *Request) BuildURL() (string, error) {
	if len(r.Query) == 0 {
		return r.URL, nil
	}

	u, err := url.Parse(r.URL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for k, vs := range r.Query {
		for _, v := range vs {
			q.Add(k, v)
		}
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}
