package json

import "github.com/ydnar/codec"

type AddressBook struct {
	Addresses []Address
}

func (b *AddressBook) Marshal(enc codec.Encoder) error {
	return codec.EncodeStruct(enc, b)
}

func (b *AddressBook) MarshalStruct(enc codec.StructEncoder) error {
	return enc.Encode("addresses", codec.Slice(&b.Addresses))
}

type Address struct {
	Name       string
	Number     int
	Street     string
	City       *City
	PostalCode string
}

func (a *Address) MarshalStruct(enc codec.StructEncoder) error {
	enc.Encode("name", a.Name)
	enc.Encode("number", a.Number)
	enc.Encode("street", a.Street)
	return enc.Encode("city", a.City)
}

type City struct {
	Name  string
	State *State
}

func (c *City) MarshalStruct(enc codec.StructEncoder) error {
	enc.Encode("name", c.Name)
	return enc.Encode("state", c.State)
}

type State struct {
	Name string
	Code string
}

func (c *State) MarshalStruct(enc codec.StructEncoder) error {
	enc.Encode("name", c.Name)
	return enc.Encode("code", c.Code)
}
