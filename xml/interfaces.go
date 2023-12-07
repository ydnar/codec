package xml

import "github.com/ydnar/codec"

type Resolver interface {
	codec.Resolver
	ResolveXMLCodec(Name) codec.Codec
}

type Name struct {
	ns   *NS
	name string
}

type NS struct {
	uri    string
	prefix string
}

func (ns NS) Equal(other NS) bool {
	return ns.uri == other.uri
}
