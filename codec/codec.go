package codec

type Codec interface {
	ContentType() string
	Encode(v any) ([]byte, error)
	Decode(data []byte, v any) error
}
