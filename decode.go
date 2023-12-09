package codec

import (
	"encoding"
	"strconv"
	"unsafe"
)

// DecodeNil calls DecodeNil on v if v implements NilDecoder.
// Returns true if v implements NilDecoder and a decode was attempted.
func DecodeNil(v any) (bool, error) {
	if v, ok := v.(NilDecoder); ok {
		return true, v.DecodeNil()
	}
	return false, nil
}

// DecodeValue decodes a string value into v by calling
// DecodeString, DecodeBytes, DecodeBoolString, and DecodeNumber,
// returning after the first attempted decode.
// Returns false, nil if unable to decode into v.
func DecodeValue(v any, val string) (bool, error) {
	if ok, err := DecodeString(v, val); ok {
		return ok, err
	}
	if ok, err := DecodeBytes(v, []byte(val)); ok {
		return ok, err
	}
	if ok, err := DecodeBoolString(v, val); ok {
		return ok, err
	}
	if ok, err := DecodeNumber(v, val); ok {
		return ok, err
	}
	return false, nil
}

// DecodeBool decodes a boolean value into v.
// If *v is a pointer to a bool, then a bool will be allocated.
// If v implements BoolDecoder, then DecodeBool(b) is called.
// Returns true if v matches a known type and a decode was attempted.
func DecodeBool(v any, b bool) (bool, error) {
	switch v := v.(type) {
	case *bool:
		*v = b
		return true, nil
	case **bool:
		*Must(v) = b
		return true, nil
	case BoolDecoder:
		return true, v.DecodeBool(b)
	}
	return false, nil
}

// DecodeBoolString decodes a string representing a boolean value into v.
//
// The values "", "0", "false", and "FALSE" are considered false.
// The values "1", "true", and "TRUE" are considered true.
// All other values are ignored and (false, nil) will be returned.
//
// If *v is a pointer to a bool, then a bool will be allocated.
// If v implements BoolDecoder, then DecodeBool(b) is called.
//
// Returns true if a decode was attempted.
func DecodeBoolString(v any, s string) (bool, error) {
	if s == "" || s == "0" || s == "false" || s == "FALSE" {
		return DecodeBool(v, false)
	}
	if s == "1" || s == "true" || s == "TRUE" {
		return DecodeBool(v, true)
	}
	return false, nil
}

// DecodeNumber decodes a number encoded as a string into v.
// The following core types are supported:
// int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, and float64.
// Pointers to the above types are also supported, and will be allocated if necessary.
// The interface types IntDecoder, and FloatDecoder are also supported.
// Returns true if v matches a known type and a decode was attempted.
func DecodeNumber(v any, n string) (bool, error) {
	switch v := v.(type) {
	case *int:
		return decodeSignedValue(v, n)
	case **int:
		return decodeSignedValue(Must(v), n)
	case *int8:
		return decodeSignedValue(v, n)
	case **int8:
		return decodeSignedValue(Must(v), n)
	case *int16:
		return decodeSignedValue(v, n)
	case **int16:
		return decodeSignedValue(Must(v), n)
	case *int32:
		return decodeSignedValue(v, n)
	case **int32:
		return decodeSignedValue(Must(v), n)
	case *int64:
		return decodeSignedValue(v, n)
	case **int64:
		return decodeSignedValue(Must(v), n)

	case *uint:
		return decodeUnsignedValue(v, n)
	case **uint:
		return decodeUnsignedValue(Must(v), n)
	case *uint8:
		return decodeUnsignedValue(v, n)
	case **uint8:
		return decodeUnsignedValue(Must(v), n)
	case *uint16:
		return decodeUnsignedValue(v, n)
	case **uint16:
		return decodeUnsignedValue(Must(v), n)
	case *uint32:
		return decodeUnsignedValue(v, n)
	case **uint32:
		return decodeUnsignedValue(Must(v), n)
	case *uint64:
		return decodeUnsignedValue(v, n)
	case **uint64:
		return decodeUnsignedValue(Must(v), n)

	case *float32:
		return decodeFloatValue(v, n)
	case **float32:
		return decodeFloatValue(Must(v), n)
	case *float64:
		return decodeFloatValue(v, n)
	case **float64:
		return decodeFloatValue(Must(v), n)

	case IntDecoder[int]:
		return decodeSigned(v, n)
	case IntDecoder[int8]:
		return decodeSigned(v, n)
	case IntDecoder[int16]:
		return decodeSigned(v, n)
	case IntDecoder[int32]:
		return decodeSigned(v, n)
	case IntDecoder[int64]:
		return decodeSigned(v, n)

	case IntDecoder[uint]:
		return decodeUnsigned(v, n)
	case IntDecoder[uint8]:
		return decodeUnsigned(v, n)
	case IntDecoder[uint16]:
		return decodeUnsigned(v, n)
	case IntDecoder[uint32]:
		return decodeUnsigned(v, n)
	case IntDecoder[uint64]:
		return decodeUnsigned(v, n)

	case FloatDecoder[float32]:
		return decodeFloat(v, n)
	case FloatDecoder[float64]:
		return decodeFloat(v, n)
	}

	return false, nil
}

func decodeSignedValue[T Signed](v *T, n string) (bool, error) {
	i, err := strconv.ParseInt(n, 10, int(unsafe.Sizeof(*v)))
	if err != nil {
		return true, err
	}
	*v = T(i)
	return true, nil
}

func decodeSigned[T Signed](v IntDecoder[T], n string) (bool, error) {
	var x T
	i, err := strconv.ParseInt(n, 10, int(unsafe.Sizeof(x))*8)
	if err != nil {
		return true, err
	}
	return true, v.DecodeInt(T(i))
}

func decodeUnsignedValue[T Unsigned](v *T, n string) (bool, error) {
	i, err := strconv.ParseUint(n, 10, int(unsafe.Sizeof(*v)))
	if err != nil {
		return true, err
	}
	*v = T(i)
	return true, nil
}

func decodeUnsigned[T Unsigned](v IntDecoder[T], n string) (bool, error) {
	var x T
	i, err := strconv.ParseUint(n, 10, int(unsafe.Sizeof(x))*8)
	if err != nil {
		return true, err
	}
	return true, v.DecodeInt(T(i))
}

func decodeFloatValue[T Float](v *T, n string) (bool, error) {
	f, err := strconv.ParseFloat(n, int(unsafe.Sizeof(*v)))
	if err != nil {
		return true, err
	}
	*v = T(f)
	return true, nil
}

func decodeFloat[T Float](v FloatDecoder[T], n string) (bool, error) {
	var x T
	f, err := strconv.ParseFloat(n, int(unsafe.Sizeof(x))*8)
	if err != nil {
		return true, err
	}
	return true, v.DecodeFloat(T(f))
}

// DecodeString decodes s into v. The following types are supported:
// string, *string, and StringDecoder.
// Returns true if v matches a known type and a decode was attempted.
func DecodeString(v any, s string) (bool, error) {
	switch v := v.(type) {
	case *string:
		*v = s
		return true, nil
	case **string:
		*v = &s
		return true, nil
	case StringDecoder:
		return true, v.DecodeString(s)
	}
	return false, nil
}

// DecodeBytes decodes data into v. The following types are supported:
// []byte, BytesDecoder, encoding.BinaryUnmarshaler, and encoding.TextUnmarshaler.
func DecodeBytes(v any, data []byte) (bool, error) {
	switch v := v.(type) {
	case *[]byte:
		Resize(v, len(data))
		copy(*v, data)
		return true, nil
	case BytesDecoder:
		return true, v.DecodeBytes(data)
	case encoding.BinaryUnmarshaler:
		return true, v.UnmarshalBinary(data)
	case encoding.TextUnmarshaler:
		return true, v.UnmarshalText(data)
	}
	return false, nil
}

// DecodeText decodes text into v. The following types are supported:
// []byte, TextDecoder, and encoding.TextUnmarshaler.
func DecodeText(v any, text []byte) (bool, error) {
	switch v := v.(type) {
	case *[]byte:
		Resize(v, len(text))
		copy(*v, text)
		return true, nil
	case *string:
		*v = string(text)
		return true, nil
	case **string:
		s := string(text)
		*v = &s
		return true, nil
	case TextDecoder:
		return true, v.DecodeText(text)
	case encoding.TextUnmarshaler:
		return true, v.UnmarshalText(text)
	}
	return false, nil
}

// DecodeSlice adapts slice s into an ElementDecoder and decodes it.
func DecodeSlice[S ~[]E, E comparable](dec Decoder, s *S) error {
	return dec.Decode(Slice(s))
}

// DecodeMap adapts a string-keyed map m into a FieldDecoder and decodes it.
func DecodeMap[M ~map[K]V, K ~string, V any](dec Decoder, m *M) error {
	return dec.Decode(Map(m))
}
