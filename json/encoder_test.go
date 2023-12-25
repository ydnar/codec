package json

import "github.com/ydnar/codec"

type AddressBook struct {
	Addresses []Address
}

func (b *AddressBook) MarshalEncode(enc codec.Encoder) error {
	e := codec.EncodeStruct(enc, "AddressBook")
	return e.Encode("addresses", codec.Slice(&b.Addresses))
}

func (b *AddressBook) MarshalStruct(enc codec.StructEncoder) error {
	return enc.Encode("addresses", codec.Slice(&b.Addresses))
}

func (b *AddressBook) MarshalFields(enc codec.FieldEncoder) error {
	return enc.EncodeField("addresses", codec.Slice(&b.Addresses))
}

type Address struct {
	Name       string
	Number     int
	Street     string
	City       *City
	PostalCode string
}

func (a *Address) MarshalFields(enc codec.FieldEncoder) error {
	enc.EncodeField("name", a.Name)
	enc.EncodeField("number", a.Number)
	enc.EncodeField("street", a.Street)
	return enc.EncodeField("city", a.City)
}

type City struct {
	Name  string
	State *State
}

func (c *City) MarshalFields(enc codec.FieldEncoder) error {
	enc.EncodeField("name", c.Name)
	return enc.EncodeField("state", c.State)
}

type State struct {
	Name string
	Code string
}

func (c *State) MarshalFields(enc codec.FieldEncoder) error {
	enc.EncodeField("name", c.Name)
	return enc.EncodeField("code", c.Code)
}
