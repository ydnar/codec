package xml

import (
	"errors"

	"github.com/ydnar/codec"
)

// Document represents an XML document.
// To decode the root node of an XML document,
// assign Root and pass Document to Decoder.Decode().
// If successful, Name will contain the name of the root XML node.
type Document struct {
	Name string
	Root any
}

func (doc *Document) DecodeElement(dec codec.Decoder, i int, name string) error {
	if i > 0 {
		return ErrMultipleRootNodes
	}
	doc.Name = name
	return dec.Decode(doc.Root)
}

func (doc *Document) DecodeXMLElement(dec codec.Decoder, name Name) error {
	return dec.Decode(nil)
}

func (doc *Document) DecodeXMLAttr(dec codec.Decoder, name Name, value []byte) error {
	return dec.Decode(nil)
}

// TODO: make this a struct with useful context
var ErrMultipleRootNodes = errors.New("xml: multiple root nodes")
