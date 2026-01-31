package client

import (
	"context"

	"github.com/imattdu/httpclient/codec"
	"github.com/imattdu/httpclient/request"
	"github.com/imattdu/httpclient/response"
)

func (c *Client) Get(ctx context.Context, url string, opts ...request.Option) (*response.Response, error) {
	req := request.New("GET", url, opts...)
	return c.Do(ctx, req)
}

func (c *Client) Post(ctx context.Context, url string, body any, cc codec.Codec, opts ...request.Option) (*response.Response, error) {
	opts = append(opts, request.WithBody(body, cc))
	req := request.New("POST", url, opts...)
	return c.Do(ctx, req)
}

func (c *Client) PostJSON(ctx context.Context, url string, body any, opts ...request.Option) (*response.Response, error) {
	opts = append(opts, request.WithBody(body, codec.NewJSON()))
	req := request.New("POST", url, opts...)
	return c.Do(ctx, req)
}
