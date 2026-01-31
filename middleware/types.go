package middleware

import (
	"context"

	"httpclient/request"
	"httpclient/response"
)

type Handler func(ctx context.Context, req *request.Request) (*response.Response, error)

type Middleware func(Handler) Handler
