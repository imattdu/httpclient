package test

import (
	"context"
	"fmt"

	"github.com/imattdu/httpclient/client"
	"github.com/imattdu/httpclient/config"
	"github.com/imattdu/httpclient/middleware"
	"testing"
	"time"
)

func TestJSONCodec_EncodeDecode(t *testing.T) {
	if err := config.Init(config.WithConfigFile("/Users/matt/workspace/go/src/httpclient/test/httpclient.yaml")); err != nil {
		fmt.Println(err.Error())
	}

	timeout := time.Second * 10
	cli := client.NewDefault(
		client.WithMiddleware(middleware.Logger()),
		client.WithConfig(&config.Config{
			Timeout: &timeout,
		}))

	res, err := cli.PostJSON(context.Background(), "http://wwww1.baidu.com",
		map[string]interface{}{
			"name": "matt",
		})
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("77%#v\n", string(res.RawBody))
}
