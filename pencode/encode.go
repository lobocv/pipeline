package pencode

// Encoder is the interface for encoding output payloads into []byte in order to be written
type Encoder interface {
	Encode(v interface{}) ([]byte, error)
}

// Decoder is the interface for decoding input payloads from []byte into structures to be used in processing
type Decoder interface {
	Decode(b []byte) (interface{}, error)
}

// allocator is a function that allocates a struct
type allocator func() interface{}
