package json

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/ydnar/codec"
)

type Decoder struct {
	dec *json.Decoder
	r   []codec.Resolver
}

func NewDecoder(r io.Reader, resolvers ...codec.Resolver) *Decoder {
	dec := json.NewDecoder(r)
	dec.UseNumber()
	return &Decoder{
		dec: dec,
		r:   resolvers,
	}
}

func (dec *Decoder) Decode(v any) error {
	if c := codec.Resolve(v, dec.r...); c != nil {
		v = c
	}

	err := dec.decodeToken(v)
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (dec *Decoder) decodeToken(v any) error {
	tok, err := dec.dec.Token()
	if err != nil {
		return err
	}
	if tok == nil {
		_, err := codec.UnmarshalNil(v)
		return err
	}

	switch tok := tok.(type) {
	case bool:
		_, err := codec.UnmarshalBool(v, tok)
		return err
	case json.Number:
		if ok, err := codec.UnmarshalDecimal(v, string(tok)); ok {
			return err
		}
		if ok, err := codec.UnmarshalString(v, string(tok)); ok {
			return err
		}
	case string:
		_, err := codec.UnmarshalValue(v, string(tok))
		return err
	case json.Delim:
		switch tok {
		case '{':
			return dec.decodeObject(v)
		case '[':
			return dec.decodeArray(v)
		default:
			return fmt.Errorf("unexpected JSON token %v at offset %d", tok, dec.dec.InputOffset())
		}
	}

	return nil
}

// decodeObject decodes a JSON object into v.
// It expects that the initial { token has already been decoded.
func (dec *Decoder) decodeObject(v any) error {
	if d, ok := v.(codec.FieldUnmarshaler); ok {
		for i := 0; dec.dec.More(); i++ {
			name, err := dec.stringToken()
			if err != nil {
				return err
			}
			once := &onceDecoder{Decoder: dec}
			err = d.UnmarshalField(once, name)
			if err != nil {
				return err
			}
			if once.calls == 0 {
				err = dec.Decode(nil)
				if err != nil {
					return err
				}
			}
		}
	} else {
		for dec.dec.More() {
			err := dec.Decode(nil)
			if err != nil {
				return err
			}
		}
	}

	tok, err := dec.dec.Token()
	if err != nil {
		return err
	}
	if tok != json.Delim('}') {
		return fmt.Errorf("unexpected JSON token %v at offset %d", tok, dec.dec.InputOffset())
	}

	return nil
}

// decodeArray decodes a JSON array into v.
// It expects that the initial [ token has already been decoded.
func (dec *Decoder) decodeArray(v any) error {
	if d, ok := v.(codec.SeqUnmarshaler); ok {
		for i := 0; dec.dec.More(); i++ {
			once := &onceDecoder{Decoder: dec}
			err := d.UnmarshalSeq(once, i)
			if err != nil {
				return err
			}
			if once.calls == 0 {
				err = dec.Decode(nil)
				if err != nil {
					return err
				}
			}
		}
	} else {
		for dec.dec.More() {
			err := dec.Decode(nil)
			if err != nil {
				return err
			}
		}
	}

	tok, err := dec.dec.Token()
	if err != nil {
		return err
	}
	if tok != json.Delim(']') {
		return fmt.Errorf("unexpected JSON token %v at offset %d", tok, dec.dec.InputOffset())
	}

	return nil
}

func (dec *Decoder) stringToken() (string, error) {
	tok, err := dec.dec.Token()
	if err != nil {
		return "", err
	}
	s, ok := tok.(string)
	if !ok {
		return "", fmt.Errorf("unexpected JSON token %v at offset %d", tok, dec.dec.InputOffset())
	}
	return s, nil
}

type decodeFunc func(codec.Decoder, int, string) error

type onceDecoder struct {
	*Decoder
	calls int
}

func (dec *onceDecoder) Decode(v any) error {
	dec.calls++
	if dec.calls > 1 {
		return fmt.Errorf("unexpected call to Decode (%d > 1)", dec.calls)
	}
	return dec.Decoder.Decode(v)
}
