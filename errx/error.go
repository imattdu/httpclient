package errx

import (
	"errors"
	"fmt"
)

type Kind string

const (
	// ErrInitConfig 配置 / 初始化
	ErrInitConfig Kind = "init_config"

	// ErrCodecNotExist codec
	ErrCodecNotExist Kind = "codec_not_exist"
	ErrEncode        Kind = "encode"
	ErrDecode        Kind = "decode"

	// ErrBuildRequest http 构建 & IO
	ErrBuildRequest Kind = "build_request"
	ErrReadBody     Kind = "read_body"

	// ErrNetwork 网络 / 超时
	ErrNetwork Kind = "network"
	ErrTimeout Kind = "timeout"

	// ErrHTTP HTTP 协议错误（有 response）
	ErrHTTP Kind = "http"
)

type Error struct {
	Kind       Kind
	Err        error
	StatusCode int
	Body       []byte
}

func (e *Error) Error() string {
	if e.Err == nil {
		return string(e.Kind)
	}
	return fmt.Sprintf("%s: %v", e.Kind, e.Err)
}

// Unwrap 支持 errors.Is / errors.As
func (e *Error) Unwrap() error {
	return e.Err
}

func IsNetwork(err error) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Kind == ErrNetwork || e.Kind == ErrTimeout
	}
	return false
}

func IsTimeout(err error) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Kind == ErrTimeout
	}
	return false
}

func IsHTTP(err error) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Kind == ErrHTTP
	}
	return false
}

func IsCodec(err error) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Kind == ErrEncode || e.Kind == ErrDecode || e.Kind == ErrCodecNotExist
	}
	return false
}
