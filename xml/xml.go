package xml

import (
	"errors"

	"github.com/ydnar/codec"
)

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

// Document represents an XML document.
// To decode the root node of an XML document,
// assign Root and pass Document to Decoder.Decode().
// If successful, Name will contain the name of the root XML node.
type Document struct {
	Name Name
	Root any
}

func (doc *Document) DecodeXMLElement(dec codec.Decoder, name Name) error {
	if doc.Name != (Name{}) {
		return ErrMultipleRootNodes
	}
	doc.Name = name
	return dec.Decode(doc.Root)
}

// TODO: make this a struct with useful context
var ErrMultipleRootNodes = errors.New("xml: multiple root nodes")
