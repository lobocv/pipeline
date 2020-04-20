package pencode

type Encoder interface {
	Encode(v interface{}) ([]byte, error)
}

type Decoder interface {
	Decode(b []byte) (interface{}, error)
}

type allocator func() interface{}
