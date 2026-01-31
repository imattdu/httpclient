package transport

import (
	"context"
	"net/http"
)

type Transport interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}

// Middleware wraps Transport
type Middleware func(Transport) Transport

// TransportFunc adapts a function to Transport
type TransportFunc func(ctx context.Context, req *http.Request) (*http.Response, error)

func (f TransportFunc) Do(
	ctx context.Context,
	req *http.Request,
) (*http.Response, error) {
	return f(ctx, req)
}

// HTTPTransport is the real bottom transport
type HTTPTransport struct {
	client *http.Client
}

func New(client *http.Client) *HTTPTransport {
	return &HTTPTransport{client: client}
}

func (t *HTTPTransport) Do(
	ctx context.Context,
	req *http.Request,
) (*http.Response, error) {
	return t.client.Do(req)
}
