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

// UnmarshalValue unmarshals a string value into v.
// It calls [UnmarshalString], [UnmarshalText], [UnmarshalBoolString], and [UnmarshalDecimal],
// returning after the first matching unmarshaler.
// Returns false, nil if unable to unmarshal val into v.
func UnmarshalValue(v any, val string) (bool, error) {
	if ok, err := UnmarshalString(v, val); ok {
		return ok, err
	}
	if ok, err := UnmarshalText(v, []byte(val)); ok {
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

func unmarshalSignedDecimal[T int | int8 | int16 | int32 | int64](v *T, n string) (bool, error) {
	i, err := strconv.ParseInt(n, 10, int(unsafe.Sizeof(*v)))
	if err != nil {
		return true, err
	}
	*v = T(i)
	return true, nil
}

func unmarshalUnsignedDecimal[T uint | uint8 | uint16 | uint32 | uint64](v *T, n string) (bool, error) {
	i, err := strconv.ParseUint(n, 10, int(unsafe.Sizeof(*v)))
	if err != nil {
		return true, err
	}
	*v = T(i)
	return true, nil
}

func unmarshalFloatDecimal[T float32 | float64](v *T, n string) (bool, error) {
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

// UnmarshalString unmarshals string s into v.
// Supported types of v: *string, **string, [StringUnmarshaler], and [encoding.TextUnmarshaler].
// Returns true if v is a supported type.
func UnmarshalString(v any, s string) (bool, error) {
	switch v := v.(type) {
	case *string:
		*v = s
		return true, nil
	case **string:
		*v = &s
		return true, nil
	case StringUnmarshaler:
		return true, v.UnmarshalString(s)
	case encoding.TextUnmarshaler:
		return true, v.UnmarshalText([]byte(s))
	}
	return false, nil
}

// UnmarshalText unmarshals text into v.
// Supported types of v: *[]byte, *string, **string, or [encoding.TextUnmarshaler].
// If v is a byte slice or string, v will be set to text.
func UnmarshalText(v any, text []byte) (bool, error) {
	switch v := v.(type) {
	case *[]byte:
		Expand(v, len(text))
		copy(*v, text)
		return true, nil
	case *string:
		*v = string(text)
		return true, nil
	case **string:
		s := string(text)
		*v = &s
		return true, nil
	case encoding.TextUnmarshaler:
		return true, v.UnmarshalText(text)
	}
	return false, nil
}

// AppendText appends text onto v.
// Unlike [UnmarshalText], this will append to, rather than replace the value of v.
//
// Supported types of v: *[]byte, *string, **string, [TextAppender], and [encoding.TextUnmarshaler].
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
	case TextAppender:
		return true, v.AppendText(text)
	case encoding.TextUnmarshaler:
		return true, v.UnmarshalText(text)
	}
	return false, nil
}

// UnmarshalBinary unmarshals binary data into v.
// Supported types of v: *[]byte, or [encoding.BinaryUnmarshaler].
// Returns true if v is a supported type.
func UnmarshalBinary(v any, data []byte) (bool, error) {
	switch v := v.(type) {
	case *[]byte:
		Expand(v, len(data))
		copy(*v, data)
		return true, nil
	case encoding.BinaryUnmarshaler:
		return true, v.UnmarshalBinary(data)
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
