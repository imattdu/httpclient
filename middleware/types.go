package middleware

import (
	"context"

	"github.com/imattdu/httpclient/request"
	"github.com/imattdu/httpclient/response"
)

type Handler func(ctx context.Context, req *request.Request) (*response.Response, error)

type Middleware func(Handler) Handler
