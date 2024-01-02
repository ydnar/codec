package codec

import (
	"cmp"
	"slices"
	"unsafe"
)

// Must ensures that the pointer at *p is non-nil.
// If *p is nil, a new value of type T will be allocated.
func Must[T any](p **T) *T {
	if *p == nil {
		*p = new(T)
	}
	return *p
}

// Expand ensures the length of slice s is greater than or equal to size.
func Expand[S ~[]E, E any](s *S, size int) {
	if size > len(*s) {
		*s = append(*s, make([]E, size-len(*s))...)
	}
}

// Slice returns a [SliceCodec] for slice s.
func Slice[S ~[]E, E any](s *S) SliceCodec {
	return (*sliceCodec[E])(unsafe.Pointer(s))
}

type SliceCodec interface {
}

// sliceCodec is an implementation of [ElementDecoder] for an arbitrary slice.
type sliceCodec[E any] []E

// UnmarshalElement implements the [ElementUnmarshaler] interface,
// dynamically resizing the slice if necessary.
func (c *sliceCodec[E]) UnmarshalElement(dec Decoder, i int) error {
	var v E
	if i >= 0 && i < len(*c) {
		v = (*c)[i]
	}
	err := dec.Decode(&v)
	if err != nil {
		return err
	}
	Expand(c, i+1)
	(*c)[i] = v
	return nil
}

func (c *sliceCodec[E]) MarshalSeq(enc Encoder) error {
	for i := range *c {
		err := enc.Encode((*c)[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// Map returns an [FieldDecoder] for map m.
func Map[M ~map[K]V, K ~string, V any](m *M) FieldDecoder {
	return &mapCodec[M, K, V]{m}
}

// mapCodec is an implementation of [FieldDecoder] for an arbitrary map with string keys.
type mapCodec[M ~map[K]V, K ~string, V any] struct{ M *M }

// DecodeField implements the [FieldDecoder] interface,
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
