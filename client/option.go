package client

import (
	"github.com/imattdu/httpclient/config"
	"github.com/imattdu/httpclient/middleware"
)

type Option func(*Client)

func WithConfig(cfg *config.Config) Option {
	return func(c *Client) {
		if cfg != nil {
			c.config = cfg
		}
	}
}

func WithMiddleware(m middleware.Middleware) Option {
	return func(c *Client) {
		c.middlewares = append(c.middlewares, m)
	}
}

// WithResolver 注入配置解析器
func WithResolver(r config.Resolver) Option {
	return func(c *Client) {
		c.resolver = r
	}
}
