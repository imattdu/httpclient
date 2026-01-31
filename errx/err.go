package errx

import "fmt"

type ErrorKind string

const (
	ErrTimeout ErrorKind = "timeout"
	ErrNetwork ErrorKind = "network"
	ErrHTTP    ErrorKind = "http"
	ErrEncode  ErrorKind = "encode"
	ErrDecode  ErrorKind = "decode"
	ErrConfig  ErrorKind = "config"
	ErrCodec   ErrorKind = "codec"
)

type Error struct {
	Kind       ErrorKind
	StatusCode int
	Body       []byte
	Err        error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return fmt.Sprintf("httpclient error: %s", e.Kind)
}

func (e *Error) Unwrap() error {
	return e.Err
}
