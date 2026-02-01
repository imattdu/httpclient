package errx

type Option func(*Error)

func WithStatusCode(code int) Option {
	return func(e *Error) {
		e.StatusCode = code
	}
}

func WithBody(body []byte) Option {
	return func(e *Error) {
		e.Body = body
	}
}
