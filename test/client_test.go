package test

import (
	"context"
	"fmt"
	"httpclient/client"
	"httpclient/codec"
	"httpclient/request"
	"httpclient/transport"
	"net/http"
	"testing"
)

func TestJSONCodec_EncodeDecode(t *testing.T) {
	cli := client.New(
		transport.New(&http.Client{}),
	)

	type A struct {
		Name string `json:"name"`
	}

	resp, err := cli.Post(context.Background(), "https://1u71583ze0m4.mock.hoppscotch.io/ping/users",
		map[string]string{
			"foo": "bar",
		}, request.WithCodec(codec.JSON))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var a A
	_ = resp.Decode(&a)
	fmt.Println(resp, a)

}
