package transport

import (
	"context"
	"net/http"
)

type Transport interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}
