package xml

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/ydnar/codec"
)

type Decoder struct {
	dec *xml.Decoder
	r   []codec.Resolver

	// Current XML element
	e xml.StartElement

	// Current XML attribute to be decoded
	attr xml.Attr
}

func NewDecoder(r io.Reader, resolvers ...codec.Resolver) *Decoder {
	dec := xml.NewDecoder(r)
	return &Decoder{
		dec: dec,
		r:   resolvers,
	}
}

func (dec *Decoder) Decode(v any) error {
	if c := codec.Resolve(v, dec.r...); c != nil {
		v = c
	}

	err := dec.decode(v)
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (dec *Decoder) decode(v any) error {
	// Decode single attribute, if set.
	if dec.attr != (xml.Attr{}) {
		_, err := codec.UnmarshalValue(v, dec.attr.Value)
		return err
	}

	// Decode start element attributes, if set.
	if len(dec.e.Attr) != 0 {
		if d, ok := v.(codec.FieldUnmarshaler); ok {
			for i := 0; i < len(dec.e.Attr); i++ {
				dec.attr = dec.e.Attr[i]
				err := d.UnmarshalField(dec, flatten(dec.attr.Name))
				if err != nil {
					return err
				}
			}
		}
		dec.attr = xml.Attr{}
		dec.e.Attr = nil
	}

	// Decode child nodes.
	for {
		tok, err := dec.dec.Token()
		if err != nil {
			return err
		}

		switch tok := tok.(type) {
		// TODO: handle PIs, chardata, CDATA, etc.
		case xml.CharData:
			_, err := codec.AppendText(v, tok)
			if err != nil {
				return err
			}
		case xml.StartElement:
			err := dec.decodeElement(v, tok)
			if err != nil {
				return err
			}
		case xml.EndElement:
			if dec.e.Name != tok.Name {
				line, col := dec.dec.InputPos()
				return fmt.Errorf("mismatched end tag %q != %q at line %d, column %d",
					flatten(tok.Name), flatten(dec.e.Name), line, col)
			}
			return nil
		}
	}
}

func (dec *Decoder) decodeElement(v any, start xml.StartElement) error {
	saved := dec.e
	defer func() { dec.e = saved }()
	dec.e = start

	switch v := v.(type) {
	case codec.FieldUnmarshaler:
		once := &onceDecoder{Decoder: dec}
		err := v.UnmarshalField(once, flatten(start.Name))
		if err != nil {
			return err
		}
		if once.calls == 0 {
			return dec.Decode(nil)
		}
	case codec.TextAppender, *[]byte, *string:
		return dec.Decode(v)
	default:
		return dec.Decode(nil)
	}

	return nil
}

func flatten(name xml.Name) string {
	if name.Space == "" {
		return name.Local
	}
	return name.Space + " " + name.Local
}

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
