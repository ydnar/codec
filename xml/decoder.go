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

	// Current parent XML element
	e xml.StartElement

	// Number of decoded elements in current parent element or XML document
	count int

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
				err := d.DecodeField(dec, i, flatten(dec.attr.Name))
				if err != nil {
					return err
				}
			}
		}
		dec.e.Attr = nil
		dec.attr = xml.Attr{}
	}

	// Decode child nodes.
	for {
		tok, err := dec.dec.Token()
		if err != nil {
			return err
		}

		switch tok := tok.(type) {
		// TODO: handle PIs, chardata, CDATA, etc.
		case xml.StartElement:
			return dec.decodeElement(v, tok)
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
	count := dec.count

	// Temporarily set state
	dec.e = start
	dec.count = 0

	var err error
	switch v := v.(type) {
	case codec.ElementDecoder:
		err = v.DecodeElement(dec, count, flatten(start.Name))
	case codec.FieldDecoder:
		err = v.DecodeField(dec, count, flatten(start.Name))
	}

	// TODO: check call count

	// Restore state
	dec.e = saved
	dec.count = count + 1

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

type ignore struct{}

func (ig *ignore) DecodeElement(dec codec.Decoder, i int, name string) error {
	return nil
}

func (ig *ignore) DecodeField(dec codec.Decoder, i int, name string) error {
	return nil
}
