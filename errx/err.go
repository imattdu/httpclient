package errx

import "fmt"

type ErrorKind string

const (
	ErrInitConfig ErrorKind = "initConfig"

	ErrTimeout  ErrorKind = "timeout"
	ErrNetwork  ErrorKind = "network"
	ErrHTTP     ErrorKind = "http"
	ErrReadBody ErrorKind = "readBody"

	ErrCodecNotExist ErrorKind = "codecNotExist"

	ErrBuildRequest ErrorKind = "buildRequest"
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
	return fmt.Sprintf("kind=: %s,error: %s", e.Kind, e.Err.Error())
}

func (e *Error) Unwrap() error {
	return e.Err
}
