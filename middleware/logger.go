package middleware

import (
	"context"
	"httpclient/transport"
	"log"
	"net/http"
	"time"
)

func Logger() transport.Middleware {
	return func(next transport.Transport) transport.Transport {
		return transport.TransportFunc(
			func(ctx context.Context, req *http.Request) (*http.Response, error) {
				start := time.Now()
				resp, err := next.Do(ctx, req)
				log.Println(req.Method, req.URL, time.Since(start), err)
				return resp, err
			},
		)
	}
}
