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

// NilDecoder is the interface implemented by types that can decode from nil.
type NilDecoder interface {
	DecodeNil() error
}

// BoolDecoder is the interface implemented by types that can decode from a bool.
type BoolDecoder interface {
	DecodeBool(bool) error
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
//
// The ordinal index and name (if any) are supplied to DecodeElement
// via the int and string arguments, respectively. The string value
// may be empty if the source format (such as a JSON array) does not
// have a natural name.
type ElementDecoder interface {
	DecodeElement(Decoder, int, string) error
}

// FieldDecoder is the interface implemented by types that can decode
// fields or attributes, such as structs or string-keyed maps.
//
// The ordinal index and name (if any) of the decoded field are supplied
// to DecodeField via the int and string arguments, respectively. The string value
// may be empty if the source format (such as a JSON array) does not
// have a natural name.
type FieldDecoder interface {
	DecodeField(Decoder, int, string) error
}
