package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/imattdu/httpclient/request"
	"github.com/imattdu/httpclient/response"
)

type RetryConfig struct {
	MaxRetries int
	Interval   time.Duration

	// 是否允许重试
	RetryIf func(
		ctx context.Context,
		req *request.Request,
		resp *response.Response,
		err error,
	) bool
}

func Retry(cfg RetryConfig) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req *request.Request) (*response.Response, error) {

			var lastResp *response.Response
			var lastErr error

			for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
				fmt.Println(time.Now(), attempt, string(req.RawBody))
				// 第一次不 sleep
				if attempt > 0 {
					select {
					case <-time.After(cfg.Interval):
					case <-ctx.Done():
						return nil, ctx.Err()
					}
				}

				resp, err := next(ctx, req)
				lastResp = resp
				lastErr = err

				if !cfg.RetryIf(ctx, req, resp, err) {
					return resp, err
				}
			}

			return lastResp, lastErr
		}
	}
}
