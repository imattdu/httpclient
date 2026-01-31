package httpclient

import "errors"

type ErrorKind string

const (
	//ErrNetwork ErrorKind = "network"
	ErrTimeout ErrorKind = "timeout"
	ErrHTTP    ErrorKind = "http"
	ErrDecode  ErrorKind = "decode"
)

var (
	ErrNoCodec    = errors.New("no codec provided")
	ErrEncode     = errors.New("encode error")
	ErrNetwork    = errors.New("network error")
	ErrHTTPStatus = errors.New("http status error")
)

//type Error struct {
//	Kind       ErrorKind
//	StatusCode int
//	Body       []byte
//	Err        error
//}
//
//func (e *Error) Error() string {
//	return e.Err.Error()
//}
//
//func (e *Error) Unwrap() error {
//	return e.Err
//}
//
//var ErrNoCodec = errors.New("no codec specified")
