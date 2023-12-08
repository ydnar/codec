package xml

import (
	"strings"
	"testing"

	"github.com/ydnar/codec"
)

func TestDecoderSimple(t *testing.T) {
	x := `<simple age="1" name="hello" here="true" />`
	want := Simple{1, "hello", true}
	var v Simple
	dec := NewDecoder(strings.NewReader(x))
	err := dec.Decode(&v)
	if err != nil {
		t.Error(err)
	}
	if v != want {
		t.Errorf("Decode: got %v, expected %v", v, want)
	}
}

type Simple struct {
	Age  int
	Name string
	Here bool
}

func (s *Simple) DecodeField(dec codec.Decoder, name string) error {
	switch name {
	case "simple":
		return dec.Decode(s)
	case "age":
		return dec.Decode(&s.Age)
	case "name":
		return dec.Decode(&s.Name)
	case "here":
		return dec.Decode(&s.Here)
	}
	return nil
}

func TestDecoderComplex(t *testing.T) {
	x := `<complex length="99"><simple age="1" name="hello" here="true" /></complex>`
	want := Complex{99, Simple{Age: 1, Name: "hello", Here: true}}
	var v Complex
	dec := NewDecoder(strings.NewReader(x))
	err := dec.Decode(&v)
	if err != nil {
		t.Error(err)
	}
	if v != want {
		t.Errorf("Decode: got %v, expected %v", v, want)
	}
}

type Complex struct {
	Length int
	Simple Simple
}

func (c *Complex) DecodeField(dec codec.Decoder, name string) error {
	switch name {
	case "complex":
		return dec.Decode(c)
	case "length":
		return dec.Decode(&c.Length)
	case "simple":
		return dec.Decode(&c.Simple)
	}
	return nil
}

func TestDecoderText(t *testing.T) {
	x := `<text>Here is some text that ignores the <b>bold</b> tag.</text>`
	want := `Here is some text that ignores the bold tag.`
	var v Text
	dec := NewDecoder(strings.NewReader(x))
	err := dec.Decode(&v)
	if err != nil {
		t.Error(err)
	}
	if string(v) != want {
		t.Errorf("Decode: got %q, expected %q", string(v), want)
	}
}

type Text []byte

func (t *Text) DecodeText(text []byte) error {
	*t = append(*t, text...)
	return nil
}

func (t *Text) DecodeField(dec codec.Decoder, name string) error {
	return dec.Decode(t) // passthrough
}
