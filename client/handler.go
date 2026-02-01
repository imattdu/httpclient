package client

import (
	"context"
	"fmt"
	"github.com/imattdu/httpclient"
	"github.com/imattdu/httpclient/request"
	"github.com/imattdu/httpclient/response"
	"io"
	"net/http"
)

func newBaseHandler(c *Client) func(context.Context, *request.Request) (*response.Response, error) {
	return func(ctx context.Context, req *request.Request) (*response.Response, error) {
		// ===== 1️⃣ resolve config =====
		cfg := c.resolveConfig(req)

		if cfg.Timeout > 0 {
			fmt.Println(cfg.Timeout)
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, cfg.Timeout)
			defer cancel()
		}

		var body io.Reader
		if req.Body != nil {
			if req.Codec == nil {
				return nil, httpclient.ErrNoCodec
			}
			pr, pw := io.Pipe()
			go func() {
				defer pw.Close()
				_ = req.Codec.Encode(pw, req.Body)
			}()
			body = pr
		}

		httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL, body)
		if err != nil {
			return nil, err
		}

		httpReq.Header = req.Header.Clone()

		resp, err := c.transport.Do(ctx, httpReq)
		if err != nil {
			return nil, err
		}
		if cfg.DefaultStream {
			return &response.Response{
				StatusCode: resp.StatusCode,
				Header:     resp.Header.Clone(),
				Codec:      req.Codec,
			}, nil
		}

		defer resp.Body.Close()

		raw, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return &response.Response{
			StatusCode: resp.StatusCode,
			Header:     resp.Header.Clone(),
			RawBody:    raw,
			Codec:      req.Codec,
		}, nil
	}
}
