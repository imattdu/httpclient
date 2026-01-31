package errx

import "errors"

var (
	ErrNoCodecProvided = errors.New("no codec provided")
	ErrHTTPStatus      = errors.New("http status error")
)

func NewConfigError(err error) *Error {
	return &Error{
		Kind: ErrConfig,
		Err:  err,
	}
}

func NewEncodeError(err error) *Error {
	return &Error{
		Kind: ErrEncode,
		Err:  err,
	}
}

func NewDecodeError(err error) *Error {
	return &Error{
		Kind: ErrDecode,
		Err:  err,
	}
}

func NewCodecNotExistError(err error) *Error {
	return &Error{
		Kind: ErrCodec,
		Err:  err,
	}
}

func NewNetworkError(err error) *Error {
	return &Error{
		Kind: ErrNetwork,
		Err:  err,
	}
}

func NewTimeoutError(err error) *Error {
	return &Error{
		Kind: ErrTimeout,
		Err:  err,
	}
}

func NewHTTPError(status int, body []byte, err error) *Error {
	return &Error{
		Kind:       ErrHTTP,
		StatusCode: status,
		Body:       body,
		Err:        err,
	}
}
