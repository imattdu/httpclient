package codec

import "io"

type Codec interface {
	ContentType() string

	Encode(w io.Writer, v any) error
	Decode(r io.Reader, v any) error
}
