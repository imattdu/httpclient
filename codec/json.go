package codec

import (
	"encoding/json"
	"io"
)

type JSON struct{}

func NewJSON() *JSON {
	return &JSON{}
}

func (c *JSON) ContentType() string {
	return "application/json"
}

func (c *JSON) Encode(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	return enc.Encode(v)
}

func (c *JSON) Decode(r io.Reader, v any) error {
	dec := json.NewDecoder(r)
	return dec.Decode(v)
}
