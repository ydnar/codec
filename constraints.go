package codec

// Signed is the set of signed integer types supported by this package.
type Signed interface {
	int | int8 | int16 | int32 | int64
}

// Unsigned is the set of unsigned integer types supported by this package.
type Unsigned interface {
	uint | uint8 | uint16 | uint32 | uint64
}

// Integer is the set of integer types supported by this package.
type Integer interface {
	Signed | Unsigned
}

// Float is the set of floating-point types supported by this package.
type Float interface {
	float32 | float64
}

// Scalar is the set of types supported by [ScalarMarshaler] and [ScalarUnmarshaler].
// See https://github.com/golang/go/issues/56235 for more information.
type Scalar interface {
	bool | int64 | uint64 | float64 | complex128
}
