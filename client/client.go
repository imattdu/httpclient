package client

import (
	"context"
	"github.com/imattdu/httpclient/config"
	"github.com/imattdu/httpclient/request"
	"github.com/imattdu/httpclient/response"
	"io"
	"net/http"

	"github.com/imattdu/httpclient/middleware"
	"github.com/imattdu/httpclient/transport"
)

type Client struct {
	handler middleware.Handler

	transport   transport.Transport
	middlewares []middleware.Middleware
	config      *config.Config
	resolver    config.Resolver
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

type LimitWriter struct {
	W     io.Writer
	N     int64
	wrote int64
}

func (l *LimitWriter) Write(p []byte) (int, error) {
	if l.wrote >= l.N {
		return len(p), nil
	}

	remain := l.N - l.wrote
	if int64(len(p)) > remain {
		p = p[:remain]
	}

	n, err := l.W.Write(p)
	l.wrote += int64(n)

	return len(p), err
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

func (c *Client) Do(ctx context.Context, req *request.Request) (*response.Response, error) {
	return c.handler(ctx, req)
}
