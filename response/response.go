package response

import (
	"httpclient/codec"
	"net/http"
)

type Response struct {
	StatusCode int
	Header     http.Header
	RawBody    []byte
	Codec      codec.Codec
}

func (r *Response) Decode(v any) error {
	if r.Codec == nil {
		return nil
	}
	return r.Codec.Decode(r.RawBody, v)
}
