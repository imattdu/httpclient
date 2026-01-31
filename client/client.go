package client

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"httpclient"
	"httpclient/request"
	"httpclient/response"
	"httpclient/transport"
)

type Client struct {
	transport transport.Transport
}

func New(
	base transport.Transport,
	mw ...transport.Middleware,
) *Client {

	t := base
	for i := len(mw) - 1; i >= 0; i-- {
		t = mw[i](t)
	}

	return &Client{transport: t}
}

func NewDefault(
	mw ...transport.Middleware,
) *Client {
	return New(
		transport.New(http.DefaultClient),
		mw...,
	)
}

func (c *Client) Do(
	ctx context.Context,
	req *request.Request,
) (*response.Response, error) {

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

	httpResp, err := c.transport.Do(ctx, httpReq)
	if err != nil {
		return nil, httpclient.ErrNetwork
	}
	defer func() {
		_ = httpResp.Body.Close()
	}()

	raw, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	resp := &response.Response{
		StatusCode: httpResp.StatusCode,
		Header:     httpResp.Header.Clone(),
		RawBody:    raw,
		Codec:      req.Codec,
	}

	if httpResp.StatusCode >= 400 {
		return resp, httpclient.ErrHTTPStatus
	}

	return resp, nil
}
