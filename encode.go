package codec

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
