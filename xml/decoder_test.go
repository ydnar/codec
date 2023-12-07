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
		t.Errorf("Decode: got %v, expdected %v", v, want)
	}
}

type Simple struct {
	Age  int
	Name string
	Here bool
}

func (s *Simple) DecodeElement(dec codec.Decoder, i int, name string) error {
	// fmt.Printf("DecodeElement(dec, %d, %q)\n", i, name)
	return dec.Decode(s)
}

func (s *Simple) DecodeField(dec codec.Decoder, i int, name string) error {
	// fmt.Printf("DecodeField(dec, %d, %q)\n", i, name)
	switch name {
	case "age":
		return dec.Decode(&s.Age)
	case "name":
		return dec.Decode(&s.Name)
	case "here":
		return dec.Decode(&s.Here)
	}
	return nil
}
