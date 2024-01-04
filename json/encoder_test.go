package json

import "github.com/ydnar/codec"

type AddressBook struct {
	Addresses []Address
}

func (b *AddressBook) Marshal(enc codec.Encoder) error {
	e := codec.EncodeStruct(enc)
	return e.Encode("addresses", codec.Slice(&b.Addresses))
}

type Address struct {
	Name       string
	Number     int
	Street     string
	City       *City
	PostalCode string
}

func (a *Address) Marshal(enc codec.Encoder) error {
	e := codec.EncodeStruct(enc)
	e.Encode("name", a.Name)
	e.Encode("number", a.Number)
	e.Encode("street", a.Street)
	return e.Encode("city", a.City)
}

type City struct {
	Name  string
	State *State
}

func (c *City) Marshal(enc codec.Encoder) error {
	e := codec.EncodeStruct(enc)
	e.Encode("name", c.Name)
	return e.Encode("state", c.State)
}

type State struct {
	Name string
	Code string
}

func (c *State) Marshal(enc codec.Encoder) error {
	e := codec.EncodeStruct(enc)
	e.Encode("name", c.Name)
	return e.Encode("code", c.Code)
}
