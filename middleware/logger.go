package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/imattdu/httpclient/request"
	"github.com/imattdu/httpclient/response"
)

func Logger() Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req *request.Request) (*response.Response, error) {
			fmt.Println("-log-开始请求--")
			start := time.Now()
			resp, err := next(ctx, req)

			logMap := map[string]interface{}{
				"req_body": string(req.RawBody),
				"cost":     time.Since(start).Milliseconds(),
			}
			msg, _ := json.Marshal(logMap)
			fmt.Println("-log-请求结束--||", string(msg))

			return resp, err
		}
	}
}
