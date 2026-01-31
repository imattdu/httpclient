package response

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"httpclient/codec"
)

type Response struct {
	StatusCode int
	Header     http.Header
	RawBody    []byte
	Stream     io.ReadCloser // stream 响应
	Codec      codec.Codec

	HTTPResponse *http.Response
}

func (r *Response) Decode(v any) error {
	if r.Codec == nil {
		return nil
	}

	// ❌ stream 响应不允许在这里 decode
	if r.Stream != nil {
		return errors.New("cannot Decode stream response; consume Stream directly")
	}

	if r.RawBody == nil {
		return nil
	}

	return r.Codec.Decode(bytes.NewReader(r.RawBody), v)
}
