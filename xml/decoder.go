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
		_, err := codec.DecodeValue(v, dec.attr.Value)
		return err
	}

	// Decode start element attributes, if set.
	if len(dec.e.Attr) != 0 {
		if d, ok := v.(codec.FieldDecoder); ok {
			for i := 0; i < len(dec.e.Attr); i++ {
				dec.attr = dec.e.Attr[i]
				err := d.DecodeField(dec, flatten(dec.attr.Name))
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
	// Save state
	saved := dec.e

	// Temporarily set state
	dec.e = start

	once := &onceDecoder{Decoder: dec}
	var err error
	switch v := v.(type) {
	case codec.FieldDecoder:
		err = v.DecodeField(once, flatten(start.Name))
	}

	// TODO: check Decode call count and call err = dec.Decode(nil)
	if once.calls == 0 {
		err = dec.Decode(nil)
	}

	// Restore state
	dec.e = saved

	return err
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
