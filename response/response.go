package response

import (
	"bytes"
	"errors"
	"github.com/imattdu/httpclient/errx"
	"net/http"

	"github.com/imattdu/httpclient/codec"
)

type Response struct {
	StatusCode int
	Header     http.Header
	RawBody    []byte
	Stream     bool
	Codec      codec.Codec

	HTTPResponse *http.Response
}

func (r *Response) Decode(v any) error {
	if r.Codec == nil {
		return nil
	}

	if r.Stream {
		return errx.NewDecodeError(errors.New("cannot Decode stream response; consume Stream directly"))
	}

	if r.RawBody == nil {
		return nil
	}

	return r.Codec.Decode(bytes.NewReader(r.RawBody), v)
}
