package codec

func Encode(enc Encoder, v any) error {
	switch v := v.(type) {
	case nil:
		return EncodeNil(enc)

	case Marshaler:
		return v.Marshal(enc)

	case ScalarMarshaler[bool]:
		b, err := v.MarshalScalar()
		if err != nil {
			return nil
		}
		return EncodeBool(enc, b)

	case ScalarMarshaler[int64]:
		i, err := v.MarshalScalar()
		if err != nil {
			return nil
		}
		return EncodeInt64(enc, i)

	case ScalarMarshaler[uint64]:
		i, err := v.MarshalScalar()
		if err != nil {
			return nil
		}
		return EncodeUint64(enc, i)

	case ScalarMarshaler[float64]:
		f, err := v.MarshalScalar()
		if err != nil {
			return nil
		}
		return EncodeFloat64(enc, f)

	case ScalarMarshaler[complex128]:
		c, err := v.MarshalScalar()
		if err != nil {
			return nil
		}
		return EncodeComplex128(enc, c)
	}
	return ErrNotSupported
}

func EncodeNil(enc Encoder) error {
	if enc, ok := enc.(NilEncoder); ok {
		return enc.EncodeNil()
	}
	return ErrNotSupported
}

func EncodeScalar[T Scalar](enc Encoder, v T) error {
	switch v := any(v).(type) {
	case bool:
		return EncodeBool(enc, v)
	case int64:
		return EncodeInt64(enc, v)
	case uint64:
		return EncodeUint64(enc, v)
	case float64:
		return EncodeFloat64(enc, v)
	case complex128:
		return EncodeComplex128(enc, v)
	}
	return ErrNotSupported
}

func EncodeBool(enc Encoder, b bool) error {
	if enc, ok := enc.(BoolEncoder); ok {
		return enc.EncodeBool(b)
	}
	return ErrNotSupported
}

func EncodeInt64(enc Encoder, i int64) error {
	if enc, ok := enc.(Int64Encoder); ok {
		return enc.EncodeInt64(i)
	}
	return ErrNotSupported
}

func EncodeUint64(enc Encoder, i uint64) error {
	if enc, ok := enc.(Uint64Encoder); ok {
		return enc.EncodeUint64(i)
	}
	return ErrNotSupported
}

func EncodeFloat64(enc Encoder, f float64) error {
	if enc, ok := enc.(Float64Encoder); ok {
		return enc.EncodeFloat64(f)
	}
	return ErrNotSupported
}

func EncodeComplex128(enc Encoder, f complex128) error {
	if enc, ok := enc.(Complex128Encoder); ok {
		return enc.EncodeComplex128(f)
	}
	return ErrNotSupported
}

func EncodeStruct(enc Encoder) FieldEncoder {
	switch enc := enc.(type) {
	case StructEncoder:
		return enc.EncodeStruct()
	}
	return nullFieldEncoder{}
}

func EncodeSlice(enc Encoder) Encoder {
	switch enc := enc.(type) {
	case SliceEncoder:
		return enc.EncodeSlice()
	}
	return nullEncoder{}
}

type (
	nullEncoder      struct{}
	nullFieldEncoder struct{}
)

func (nullEncoder) Encode(v any) error                   { return nil }
func (nullFieldEncoder) Encode(name string, v any) error { return nil }
