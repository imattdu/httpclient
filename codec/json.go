package codec

import (
	"encoding/json"
)

type JSONCodec struct{}

func (JSONCodec) Encode(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (JSONCodec) Decode(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (JSONCodec) ContentType() string {
	return "application/json"
}

var JSON = JSONCodec{}
