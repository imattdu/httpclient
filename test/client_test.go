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
	cli := client.NewDefault(client.WithMiddleware(middleware.Logger()), client.WithConfig(&config.Config{
		Timeout: &timeout,
	}))

	_, err := cli.Get(context.Background(), "http://5ea93d9e-355e-4572-a80d-48780f7b4397.mock.pstmn.io/go-web-v2/ping") //map[string]string{
	//	"foo": "bar",
	//})

	if err != nil {
		fmt.Println(err.Error())
	}
}
