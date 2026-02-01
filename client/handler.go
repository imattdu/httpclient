package client

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/imattdu/httpclient/errx"
	"github.com/imattdu/httpclient/request"
	"github.com/imattdu/httpclient/response"
)

func newBaseHandler(c *Client) func(context.Context, *request.Request) (*response.Response, error) {
	return func(ctx context.Context, req *request.Request) (*response.Response, error) {
		// ===== 1️⃣ resolve config =====
		cfg := c.resolveConfig(req)

		if cfg.Timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, cfg.Timeout)
			defer cancel()
		}

		result := &response.Response{
			Codec:  req.Codec,
			Stream: cfg.ResponseBodyStream,
		}
		var body io.Reader
		// 用来记录 body（只用于日志）
		var recordBuf bytes.Buffer
		if req.Body != nil {
			if req.Codec == nil {
				return result, errx.NewCodecNotExistError(errors.New("codec not exist"))
			}

			// 旁路记录（限流）
			recordWriter := &LimitWriter{
				W: &recordBuf,
				N: 8 * 1024,
			}
			if cfg.RequestBodyStream {
				pr, pw := io.Pipe()

				// Encode 写入两个地方：
				// 1. pw → HTTP client
				// 2. recordWriter → 日志旁路
				mw := io.MultiWriter(pw, recordWriter)

				go func() {
					// ⚠️ 只能在“正常完成”时 Close
					// Encode 出错必须 CloseWithError
					defer func() {
						_ = pw.Close()
					}()

					if err := req.Codec.Encode(mw, req.Body); err != nil {
						_ = pw.CloseWithError(err)
						return
					}
				}()

				// pr 只给 HTTP client
				body = pr

			} else {
				// 用于 Encode 的完整缓冲
				var buf bytes.Buffer

				// Encode 只写一次到内存
				if err := req.Codec.Encode(&buf, req.Body); err != nil {
					return result, errx.NewEncodeError(err)
				}
				// HTTP body：直接用 bytes.Reader
				body = bytes.NewReader(buf.Bytes())

				_, _ = recordWriter.Write(buf.Bytes())
				req.RawBody = recordBuf.Bytes()
			}
		}

		httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL, body)
		if err != nil {
			return result, errx.NewBuildRequestError(err)
		}

		httpReq.Header = req.Header.Clone()

		resp, err := c.transport.Do(ctx, httpReq)
		req.RawBody = recordBuf.Bytes()
		if err != nil {
			return result, errx.NewNetworkError(err)
		}
		if resp != nil {
			result.StatusCode = resp.StatusCode
			result.Header = resp.Header
		}
		if cfg.ResponseBodyStream {
			result.HTTPResponse = resp
			return result, nil
		}

		defer func() {
			_ = resp.Body.Close()
		}()

		raw, err := io.ReadAll(resp.Body)
		result.RawBody = raw
		return result, err
	}
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
