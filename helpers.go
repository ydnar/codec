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

// Slice returns a [Seq] for slice s.
// If m implements Seq, it will be returned directly.
func Slice[S ~[]E, E any](s *S) SliceCodec {
	if s, ok := any(s).(SliceCodec); ok {
		return s
	}
	return (*sliceCodec[E])(unsafe.Pointer(s))
}

type SliceCodec interface {
	Marshal(enc Encoder) error
	UnmarshalSeq(dec Decoder, i int) error
}

// sliceCodec is an implementation of [Marshaler] and [SeqUnmarshaler] for an arbitrary slice.
type sliceCodec[E any] []E

func (c *sliceCodec[E]) Marshal(enc Encoder) error {
	e := EncodeSlice(enc)
	for i := range *c {
		err := e.Encode((*c)[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// UnmarshalSeq implements the [SeqUnmarshaler] interface,
// dynamically resizing the slice if necessary.
func (c *sliceCodec[E]) UnmarshalSeq(dec Decoder, i int) error {
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

// Map returns a [MapCodec] for map m.
// If m implements MapCodec, it will be returned directly.
func Map[M ~map[K]V, K ~string, V any](m *M) MapCodec {
	if f, ok := any(m).(MapCodec); ok {
		return f
	}
	return &mapCodec[M, K, V]{m}
}

type MapCodec interface {
	Marshal(enc Encoder) error
	UnmarshalField(dec Decoder, name string) error
}

// mapCodec is an implementation of [MapCodec] for a map with string keys.
type mapCodec[M ~map[K]V, K ~string, V any] struct{ M *M }

func (c *mapCodec[M, K, V]) Marshal(enc Encoder) error {
	e := EncodeStruct(enc)
	for k, v := range *c.M {
		err := e.Encode(string(k), v)
		if err != nil {
			return err
		}
	}
	return nil
}

// UnmarshalField implements the [FieldUnmarshaler] interface,
// allocating the underlying map if necessary.
func (c *mapCodec[M, K, V]) UnmarshalField(dec Decoder, name string) error {
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
