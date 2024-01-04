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

type EncoderFunc func(v any) error

func (f EncoderFunc) Encode(v any) error {
	return f(v)
}

type ArrayEncoder interface {
	EncodeArray(int) Encoder
}

type SliceEncoder interface {
	EncodeSlice() Encoder
}

// StructEncoder is the interface implemented by Encoders that can encode a struct.
type StructEncoder interface {
	EncodeStruct() FieldEncoder
}

// FieldEncoder is the common interface implemented by a field encoder for a specific serialization format.
// A FieldEncoder is used to encode structs, or string-keyed maps.
type FieldEncoder interface {
	Encode(name string, v any) error
}

// TODO: document Marshaler
type Marshaler interface {
	Marshal(Encoder) error
}

// TODO: document Unmarshaler
type Unmarshaler interface {
	Unmarshal(Decoder) error
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

// SeqUnmarshaler is the interface implemented by types that can unmarshal
// indexed elements, such as a slice, arrays, or maps.
type SeqUnmarshaler interface {
	UnmarshalSeq(Decoder, int) error
}

// FieldUnmarshaler is the interface implemented by types that can
// unmarshal fields or attributes, such as structs or string-keyed maps.
type FieldUnmarshaler interface {
	UnmarshalField(dec Decoder, name string) error
}
