package codec

// Codec is any type that can encode or decode itself or an associated type.
type Codec any

// Resolver is the interface implemented by types that return a codec for the value at v.
// Values returned by Resolver should implement one or more encode or decode methods.
type Resolver interface {
	ResolveCodec(v any) Codec
}

// Resolve tries to resolve v with resolvers, returning the first non-nil value received.
func Resolve(v any, resolvers ...Resolver) Codec {
	for _, r := range resolvers {
		c := r.ResolveCodec(v)
		if c != nil {
			return c
		}
	}
	return nil
}

// Decoder is the common interface implemented by decoders for a specific serialization format.
type Decoder interface {
	Decode(v any) error
}

// Encoder is the common interface implemented by encoders for a specific serialization format.
type Encoder interface {
	Encode(v any) error
}

func EncodeStruct(enc Encoder, name string) StructEncoder {
	return nil
}

type StructEncoder interface {
	Encode(name string, v any) error
}

func EncodeSeq[T any](enc Encoder) SeqEncoder[T] {
	return nil
}

type SeqEncoder[T any] interface {
	Encode(v T) error
}

func EncodeMap[K comparable, V any](enc Encoder) MapEncoder[K, V] {
	return nil
}

type MapEncoder[K comparable, V any] interface {
	Encode(k K, v V) error
}
type FieldEncoder interface {
	EncodeField(name string, v any) error
}

type FieldsMarshaler interface {
	MarshalFields(enc FieldEncoder) error
}

type ElementsMarshaler interface {
	MarshalElements(enc Encoder) error
}

type Marshaler interface {
	MarshalCodec(Encoder) error
}

type Unmarshaler interface {
	UnmarshalCodec(Decoder) error
}

// ScalarMarshaler is the interface implemented by types that can marshal
// to a [Scalar] value. See https://github.com/golang/go/issues/56235 for more information.
type ScalarMarshaler[T Scalar] interface {
	MarshalScalar() (T, error)
}

// ScalarMarshaler is the interface implemented by types that can unmarshal
// from a [Scalar] value. See https://github.com/golang/go/issues/56235 for more information.
type ScalarUnmarshaler[T Scalar] interface {
	UnmarshalScalar(T) error
}

// Scalar is the set of types supported by [ScalarMarshaler] and [ScalarUnmarshaler].
// See https://github.com/golang/go/issues/56235 for more information.
type Scalar interface {
	bool | int64 | uint64 | float64 | complex128
}

// NilUnmarshaler is the interface implemented by types that can unmarshal from nil.
type NilUnmarshaler interface {
	UnmarshalNil() error
}

// BytesDecoder is the interface implemented by types that can decode from a byte slice.
// It is similar to [encoding.BinaryUnmarshaler] and [encoding.TextUnmarshaler].
type BytesDecoder interface {
	DecodeBytes([]byte) error
}

// StringUnmarshaler is the interface implemented by types that can unmarshal from a string.
type StringUnmarshaler interface {
	UnmarshalString(string) error
}

// TextAppender is the interface implemented by types that can append text data.
type TextAppender interface {
	AppendText([]byte) error
}

// ElementDecoder is the interface implemented by types that can decode
// indexed elements, such as a slice, arrays, or maps.
type ElementDecoder interface {
	DecodeElement(Decoder, int) error
}

type ElementUnmarshaler interface {
	UnmarshalElement(Decoder, int) error
}

// FieldDecoder is the interface implemented by types that can decode
// fields or attributes, such as structs or string-keyed maps.
type FieldDecoder interface {
	DecodeField(Decoder, string) error
}

// FieldUnmarshaler is the interface implemented by types that can
// unmarshal fields or attributes, such as structs or string-keyed maps.
type FieldUnmarshaler interface {
	UnmarshalField(Decoder, string) error
}

type FieldUnmarshaler2[K Key] interface {
	UnmarshalField(dec Decoder, key K) error
}

type Key interface {
	Scalar | string
}
