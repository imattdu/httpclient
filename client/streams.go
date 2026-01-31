package client

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"httpclient"
	"httpclient/request"
)

func (c *Client) DoStream(
	ctx context.Context,
	req *request.Request,
) (*http.Response, error) {

	var body io.Reader
	if req.Body != nil {
		if req.Codec == nil {
			return nil, httpclient.ErrNoCodec
		}
		b, err := req.Codec.Encode(req.Body)
		if err != nil {
			return nil, httpclient.ErrEncode
		}
		body = bytes.NewReader(b)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		req.Method,
		req.URL,
		body,
	)
	if err != nil {
		return nil, err
	}

	httpReq.Header = req.Header.Clone()

	if req.Body != nil &&
		req.Codec != nil &&
		httpReq.Header.Get("Content-Type") == "" {
		httpReq.Header.Set("Content-Type", req.Codec.ContentType())
	}

	return c.transport.Do(ctx, httpReq)
}
