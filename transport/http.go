package transport

import (
	"context"
	"net/http"
)

type HTTPTransport struct {
	client *http.Client
}

func NewHTTPTransport(client *http.Client) *HTTPTransport {
	return &HTTPTransport{client: client}
}

func (t *HTTPTransport) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	return t.client.Do(req)
}
