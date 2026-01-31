package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"httpclient/request"
	"httpclient/response"
)

func Logger() Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req *request.Request) (*response.Response, error) {
			fmt.Println("log")
			start := time.Now()
			resp, err := next(ctx, req)

			logMap := map[string]interface{}{
				"req_body": string(req.RawBody),
				"cost":     time.Since(start).Milliseconds(),
			}
			msg, _ := json.Marshal(logMap)
			fmt.Println("log", string(msg))

			return resp, err
		}
	}
}
