package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/imattdu/httpclient/client"
	"github.com/imattdu/httpclient/middleware"
	"github.com/imattdu/httpclient/request"
	"github.com/imattdu/httpclient/response"
)

func TestJSONCodec_EncodeDecode(t *testing.T) {
	cli := client.NewDefault(middleware.Logger(),
		middleware.Retry(middleware.RetryConfig{
			MaxRetries: 3,
			Interval:   time.Second,
			RetryIf: func(ctx context.Context, req *request.Request, resp *response.Response, err error) bool {
				return false
			},
		}))

	type A struct {
		Name string `json:"name"`
	}

	resp, err := cli.PostJSON(context.Background(), "https://1u71583ze0m4.mock.hoppscotch.io/ping/users",
		map[string]string{
			"foo": "bar",
		})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var a A
	_ = resp.Decode(&a)
	fmt.Println(resp, a)
}
