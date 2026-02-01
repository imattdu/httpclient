package errx

func NewInitConfigError(err error) *Error {
	return &Error{
		Kind: ErrInitConfig,
		Err:  err,
	}
}

func NewCodecNotExistError(err error) *Error {
	return &Error{
		Kind: ErrCodecNotExist,
		Err:  err,
	}
}

func NewBuildRequestError(err error) *Error {
	return &Error{
		Kind: ErrBuildRequest,
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

func NewReadBodyError(err error) *Error {
	return &Error{
		Kind: ErrReadBody,
		Err:  err,
	}
}
