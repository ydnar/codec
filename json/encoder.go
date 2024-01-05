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
	return codec.Encode(enc, v)
}

func (enc *Encoder) EncodeBool(b bool) error {
	// TODO
	return nil
}

func (enc *Encoder) EncodeInt64(i int64) error {
	// TODO
	return nil
}

func (enc *Encoder) EncodeUint64(i uint64) error {
	// TODO
	return nil
}

func (enc *Encoder) EncodeFloat64(f float64) error {
	// TODO
	return nil
}

func (enc *Encoder) EncodeString(s string) error {
	// TODO
	return nil
}

func (enc *Encoder) EncodeText(text []byte) error {
	// TODO
	return nil
}
