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

// Decoder is the interface implemented by types that can decode data into Go type(s).
type Decoder interface {
	Decode(v any) error
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

// NilUnmarshaler is the interface implemented by types that can unmarshal from nil.
type NilUnmarshaler interface {
	UnmarshalNil() error
}

// BytesDecoder is the interface implemented by types that can decode from a byte slice.
// It is similar to [encoding.BinaryUnmarshaler] and [encoding.TextUnmarshaler].
type BytesDecoder interface {
	DecodeBytes([]byte) error
}

// StringDecoder is the interface implemented by types that can decode from a string.
type StringDecoder interface {
	DecodeString(string) error
}

// TextDecoder is the interface implemented by types that can decode text from a byte slice.
type TextDecoder interface {
	DecodeText([]byte) error
}

// IntDecoder is the interface implemented by types that can decode
// from an integer value. See [Integer] for the list of supported types.
type IntDecoder[T Integer] interface {
	DecodeInt(T) error
}

// FloatDecoder is the interface implemented by types that can decode
// from a floating-point value. See [Float] for the list of supported types.
type FloatDecoder[T Float] interface {
	DecodeFloat(T) error
}

// ElementDecoder is the interface implemented by types that can decode
// indexed elements, such as a slice, arrays, or maps.
type ElementDecoder interface {
	DecodeElement(Decoder, int) error
}

// FieldDecoder is the interface implemented by types that can decode
// fields or attributes, such as structs or string-keyed maps.
type FieldDecoder interface {
	DecodeField(Decoder, string) error
}
