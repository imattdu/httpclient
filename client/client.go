package client

import (
	"context"
	"net/http"

	"github.com/imattdu/httpclient/config"
	"github.com/imattdu/httpclient/middleware"
	"github.com/imattdu/httpclient/request"
	"github.com/imattdu/httpclient/response"
	"github.com/imattdu/httpclient/transport"
)

type Client struct {
	transport   transport.Transport
	handler     middleware.Handler
	middlewares []middleware.Middleware

	config   *config.Config
	resolver config.Resolver
}

func New(t transport.Transport, opts ...Option) *Client {
	c := &Client{
		transport: t,
		config:    config.Default(),
	}

	for _, opt := range opts {
		opt(c)
	}
	if c.resolver == nil {
		c.resolver = config.DefaultManager().Resolver()
	}

	h := newBaseHandler(c)

	for i := len(c.middlewares) - 1; i >= 0; i-- {
		h = c.middlewares[i](h)
	}

	c.handler = h
	return c
}

func NewDefault(opts ...Option) *Client {
	return New(transport.NewHTTPTransport(http.DefaultClient), opts...)
}

func (c *Client) Do(ctx context.Context, req *request.Request) (*response.Response, error) {
	return c.handler(ctx, req)
}

func (c *Client) resolveConfig(req *request.Request) config.EffectiveConfig {
	if c.resolver == nil {
		return config.EffectiveConfig{}
	}

	ctx := req.Context()

	// ⭐ 把 client config 当成一层 override
	if c.config != nil {
		ctx.Override = config.Merge(c.config, ctx.Override)
	}

	return c.resolver.Resolve(ctx)
}
