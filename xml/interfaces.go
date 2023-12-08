package xml

import "github.com/ydnar/codec"

type Name struct {
	Space  string
	Prefix string
	Local  string
}

type AttrDecoder interface {
	DecodeXMLAttr(name Name, value string) error
}

type ElementDecoder interface {
	DecodeXMLElement(dec codec.Decoder, name Name) error
}
