package client

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/imattdu/httpclient"
	"github.com/imattdu/httpclient/errx"
	"github.com/imattdu/httpclient/middleware"
	"github.com/imattdu/httpclient/request"
	"github.com/imattdu/httpclient/response"
	"github.com/imattdu/httpclient/transport"
)

type Client struct {
	handler middleware.Handler
}

func New(t transport.Transport, mws ...middleware.Middleware) *Client {

	h := newBaseHandler(t)

	// 倒序包裹 middleware
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}

	return &Client{
		handler: h,
	}
}

func NewDefault(mws ...middleware.Middleware) *Client {
	return New(transport.NewHTTPTransport(http.DefaultClient), mws...)
}

func newBaseHandler(t transport.Transport) middleware.Handler {
	return func(ctx context.Context, req *request.Request) (*response.Response, error) {

		var body io.Reader
		if req.Body != nil {
			if req.Codec == nil {
				return nil, errx.NewCodecNotExistError(errors.New("codec not exist"))
			}

			pr, pw := io.Pipe()

			var logBuf bytes.Buffer
			lw := &LimitWriter{W: &logBuf, N: 4 << 10} // 最多 4KB
			tee := io.MultiWriter(pw, lw)

			done := make(chan struct{})
			go func() {
				defer close(done)
				defer func() {
					_ = pw.Close()
				}()
				if err := req.Codec.Encode(tee, req.Body); err != nil {
					_ = pw.CloseWithError(err)
				}
			}()

			body = pr
			// ⚠️ 在“请求结束后”再等 done
			go func() {
				<-done
				req.RawBody = logBuf.Bytes()
			}()
		}

		httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL, body)
		if err != nil {
			return nil, err
		}

		httpReq.Header = req.Header.Clone()

		if req.Body != nil &&
			req.Codec != nil &&
			httpReq.Header.Get("Content-Type") == "" {
			httpReq.Header.Set("Content-Type", req.Codec.ContentType())
		}

		httpResp, err := t.Do(ctx, httpReq)
		if err != nil {
			return nil, err
		}

		// ⭐⭐⭐ Stream 分支：不读 body，不 close
		if req.Stream {
			if httpResp.StatusCode >= http.StatusBadRequest {
				return &response.Response{
					StatusCode:   httpResp.StatusCode,
					Header:       httpResp.Header.Clone(),
					HTTPResponse: httpResp,
					Stream:       httpResp.Body,
				}, httpclient.ErrHTTPStatus
			}

			return &response.Response{
				StatusCode:   httpResp.StatusCode,
				Header:       httpResp.Header.Clone(),
				HTTPResponse: httpResp,
			}, nil
		}

		defer func() {
			_ = httpResp.Body.Close()
		}()

		raw, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return nil, err
		}

		resp := &response.Response{
			StatusCode: httpResp.StatusCode,
			Header:     httpResp.Header.Clone(),
			RawBody:    raw,
			Codec:      req.Codec,
		}

		if httpResp.StatusCode >= http.StatusBadRequest {
			return resp, httpclient.ErrHTTPStatus
		}

		return resp, nil
	}
}

func (c *Client) Do(ctx context.Context, req *request.Request) (*response.Response, error) {
	return c.handler(ctx, req)
}

type LimitWriter struct {
	W     io.Writer
	N     int64
	wrote int64
}

func (l *LimitWriter) Write(p []byte) (int, error) {
	if l.wrote >= l.N {
		return len(p), nil
	}

	remain := l.N - l.wrote
	if int64(len(p)) > remain {
		p = p[:remain]
	}

	n, err := l.W.Write(p)
	l.wrote += int64(n)

	return len(p), err
}
