package xml

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ydnar/codec"
)

func TestDecoder(t *testing.T) {
	x := `<foo age="1" name="hello" />`
	want := Foo{1, "hello"}
	var v Foo
	dec := NewDecoder(strings.NewReader(x))
	err := dec.Decode(&v)
	if err != nil {
		t.Error(err)
	}
	if v != want {
		t.Errorf("Decode: got %v, expdected %v", v, want)
	}
}

type Foo struct {
	Age  int
	Name string
}

func (f *Foo) DecodeElement(dec codec.Decoder, i int, name string) error {
	fmt.Printf("DecodeElement(dec, %d, %q)\n", i, name)
	return dec.Decode(f)
}

func (f *Foo) DecodeField(dec codec.Decoder, i int, name string) error {
	fmt.Printf("DecodeField(dec, %d, %q)\n", i, name)
	switch name {
	case "age":
		return dec.Decode(&f.Age)
	case "name":
		return dec.Decode(&f.Name)
	}
	return nil
}
