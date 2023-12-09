package codec

import (
	"cmp"
	"slices"
)

// Must ensures that the pointer at *p is non-nil.
// If *p is nil, a new value of type T will be allocated.
func Must[T any](p **T) *T {
	if *p == nil {
		*p = new(T)
	}
	return *p
}

// Slice returns an ElementDecoder for slice s.
func Slice[S ~[]E, E comparable](s *S) ElementDecoder {
	return &sliceCodec[S, E]{S: s}
}

// sliceCodec is an implementation of ElementDecoder for an arbitrary slice.
type sliceCodec[S ~[]E, E comparable] struct{ S *S }

// DecodeElement implements the ElementDecoder interface,
// dynamically resizing the slice if necessary.
func (c *sliceCodec[S, E]) DecodeElement(dec Decoder, i int) error {
	var v E
	if i >= 0 && i < len(*c.S) {
		v = (*c.S)[i]
	}
	err := dec.Decode(&v)
	if err != nil {
		return err
	}
	Resize(c.S, i)
	if v != (*c.S)[i] {
		(*c.S)[i] = v
	}
	return nil
}

// Resize resizes the slice s to at least len(s) == i+1,
// returning the value at s[i].
func Resize[S ~[]E, E any](s *S, i int) E {
	var e E
	if i < 0 {
		return e
	}
	if i >= len(*s) {
		*s = append(*s, make([]E, i+1-len(*s))...)
	}
	return (*s)[i]
}

// Map returns an FieldDecoder for map m.
func Map[M ~map[K]V, K ~string, V any](m *M) FieldDecoder {
	return &mapCodec[M, K, V]{m}
}

// mapCodec is an implementation of FieldDecoder for an arbitrary map with string keys.
type mapCodec[M ~map[K]V, K ~string, V any] struct{ M *M }

// DecodeField implements the FieldDecoder interface,
// allocating the underlying map if necessary.
func (c *mapCodec[M, K, V]) DecodeField(dec Decoder, name string) error {
	var v V
	err := dec.Decode(&v)
	if err != nil {
		return err
	}
	if *c.M == nil {
		*c.M = make(map[K]V)
	}
	(*c.M)[K(name)] = v
	return nil
}

// Keys returns a slice of keys for map m.
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

// SortedKeys returns a slice of keys for map m.
// Map keys must conform to cmp.Ordered.
func SortedKeys[K cmp.Ordered, V any](m map[K]V) []K {
	keys := Keys(m)
	slices.Sort(keys)
	return keys
}
