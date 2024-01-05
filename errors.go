package codec

type stringError string

func (err stringError) Error() string {
	return string(err)
}

const (
	ErrNotSupported     stringError = "codec: not supported"
	ErrBoolNotSupported stringError = "codec: bool not supported"
)
