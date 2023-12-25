package json

import (
	"encoding/json"
	"io"

	"github.com/ydnar/codec"
)

type Encoder struct {
	enc *json.Encoder
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		enc: json.NewEncoder(w),
	}
}

func (enc *Encoder) Encode(v any) error {
	switch v := v.(type) {
	case nil:
		return nil
	case codec.Marshaler:
		return v.MarshalCodec(enc)
	}
	return nil
}

// TODO
func (enc *Encoder) EncodeName(name string) (codec.Encoder, error) {
	return nil, nil
}

func (enc *Encoder) EncodeSequence() (codec.Encoder, error) {
	return nil, nil
}
