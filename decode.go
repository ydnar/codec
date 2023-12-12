package codec

import (
	"encoding"
	"strconv"
	"unsafe"
)

// UnmarshalNil calls UnmarshalNil on v if v implements NilUnmarshaler.
// Returns true if v implements NilUnmarshaler.
func UnmarshalNil(v any) (bool, error) {
	if v, ok := v.(NilUnmarshaler); ok {
		return true, v.UnmarshalNil()
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
	if ok, err := UnmarshalBoolString(v, val); ok {
		return ok, err
	}
	if ok, err := UnmarshalDecimal(v, val); ok {
		return ok, err
	}
	return false, nil
}

// UnmarshalBool unmarshals a bool into v.
// Supported types of v: *bool, **bool, and ScalarUnmarshaler[bool].
// Returns true if v matches a supported type.
func UnmarshalBool(v any, b bool) (bool, error) {
	switch v := v.(type) {
	case *bool:
		*v = b
		return true, nil
	case **bool:
		*Must(v) = b
		return true, nil
	case ScalarUnmarshaler[bool]:
		return true, v.UnmarshalScalar(b)
	}
	return false, nil
}

// UnmarshalBoolString unmarshals a string representing a bool into v.
//
// The values "", "0", "false", and "FALSE" are considered false.
// The values "1", "true", and "TRUE" are considered true.
// All other values are ignored and (false, nil) will be returned.
//
// Returns true if v matches a known type and s matches a known value.
//
// See [UnmarshalBool] for more information.
func UnmarshalBoolString(v any, s string) (bool, error) {
	if s == "" || s == "0" || s == "false" || s == "FALSE" {
		return UnmarshalBool(v, false)
	}
	if s == "1" || s == "true" || s == "TRUE" {
		return UnmarshalBool(v, true)
	}
	return false, nil
}

// UnmarshalDecimal decodes a string containing a base-10 number into v.
// The following types for v are supported:
// int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64,
// and ScalarUnmarshaler[int64 | int64 | float32 | float64].
//
// Returns true if v matches a known type.
//
// TODO: support [complex128]?
func UnmarshalDecimal(v any, n string) (bool, error) {
	switch v := v.(type) {
	case nil:
		return false, nil

	case *int:
		return unmarshalSignedDecimal(v, n)
	case **int:
		return unmarshalSignedDecimal(Must(v), n)
	case *int8:
		return unmarshalSignedDecimal(v, n)
	case **int8:
		return unmarshalSignedDecimal(Must(v), n)
	case *int16:
		return unmarshalSignedDecimal(v, n)
	case **int16:
		return unmarshalSignedDecimal(Must(v), n)
	case *int32:
		return unmarshalSignedDecimal(v, n)
	case **int32:
		return unmarshalSignedDecimal(Must(v), n)
	case *int64:
		return unmarshalSignedDecimal(v, n)
	case **int64:
		return unmarshalSignedDecimal(Must(v), n)

	case *uint:
		return unmarshalUnsignedDecimal(v, n)
	case **uint:
		return unmarshalUnsignedDecimal(Must(v), n)
	case *uint8:
		return unmarshalUnsignedDecimal(v, n)
	case **uint8:
		return unmarshalUnsignedDecimal(Must(v), n)
	case *uint16:
		return unmarshalUnsignedDecimal(v, n)
	case **uint16:
		return unmarshalUnsignedDecimal(Must(v), n)
	case *uint32:
		return unmarshalUnsignedDecimal(v, n)
	case **uint32:
		return unmarshalUnsignedDecimal(Must(v), n)
	case *uint64:
		return unmarshalUnsignedDecimal(v, n)
	case **uint64:
		return unmarshalUnsignedDecimal(Must(v), n)

	case *float32:
		return unmarshalFloatDecimal(v, n)
	case **float32:
		return unmarshalFloatDecimal(Must(v), n)
	case *float64:
		return unmarshalFloatDecimal(v, n)
	case **float64:
		return unmarshalFloatDecimal(Must(v), n)

	case ScalarUnmarshaler[int64]:
		return unmarshalScalarInt64(v, n)
	case ScalarUnmarshaler[uint64]:
		return unmarshalScalarUint64(v, n)
	case ScalarUnmarshaler[float64]:
		return unmarshalScalarFloat64(v, n)
	}

	return false, nil
}

func unmarshalSignedDecimal[T Signed](v *T, n string) (bool, error) {
	i, err := strconv.ParseInt(n, 10, int(unsafe.Sizeof(*v)))
	if err != nil {
		return true, err
	}
	*v = T(i)
	return true, nil
}

func unmarshalUnsignedDecimal[T Unsigned](v *T, n string) (bool, error) {
	i, err := strconv.ParseUint(n, 10, int(unsafe.Sizeof(*v)))
	if err != nil {
		return true, err
	}
	*v = T(i)
	return true, nil
}

func unmarshalFloatDecimal[T Float](v *T, n string) (bool, error) {
	f, err := strconv.ParseFloat(n, int(unsafe.Sizeof(*v)))
	if err != nil {
		return true, err
	}
	*v = T(f)
	return true, nil
}

func unmarshalScalarInt64(v ScalarUnmarshaler[int64], n string) (bool, error) {
	i, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		return true, err
	}
	return true, v.UnmarshalScalar(i)
}

func unmarshalScalarUint64(v ScalarUnmarshaler[uint64], n string) (bool, error) {
	i, err := strconv.ParseUint(n, 10, 64)
	if err != nil {
		return true, err
	}
	return true, v.UnmarshalScalar(i)
}

func unmarshalScalarFloat64(v ScalarUnmarshaler[float64], n string) (bool, error) {
	f, err := strconv.ParseFloat(n, 64)
	if err != nil {
		return true, err
	}
	return true, v.UnmarshalScalar(f)
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
// []byte, *string, **string, TextDecoder, and encoding.TextUnmarshaler.
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

// AppendText appends text onto v. The following types are supported:
// []byte, *string, **string, TextDecoder, and encoding.TextUnmarshaler.
func AppendText(v any, text []byte) (bool, error) {
	switch v := v.(type) {
	case *[]byte:
		*v = append(*v, text...)
		return true, nil
	case *string:
		*v += string(text)
		return true, nil
	case **string:
		*Must(v) += string(text)
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
