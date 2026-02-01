package errx

func newError(kind Kind, err error, opts ...Option) *Error {
	e := &Error{
		Kind: kind,
		Err:  err,
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func NewInitConfigError(err error) *Error {
	return newError(ErrInitConfig, err)
}

func NewCodecNotExistError(err error) *Error {
	return newError(ErrCodecNotExist, err)
}

func NewEncodeError(err error) *Error {
	return newError(ErrEncode, err)
}

func NewDecodeError(err error) *Error {
	return newError(ErrDecode, err)
}

func NewBuildRequestError(err error) *Error {
	return newError(ErrBuildRequest, err)
}

func NewNetworkError(err error) *Error {
	return newError(ErrNetwork, err)
}

func NewTimeoutError(err error) *Error {
	return newError(ErrTimeout, err)
}

func NewReadBodyError(err error) *Error {
	return newError(ErrReadBody, err)
}

func NewHTTPError(status int, body []byte, err error) *Error {
	return newError(
		ErrHTTP,
		err,
		WithStatusCode(status),
		WithBody(body),
	)
}
