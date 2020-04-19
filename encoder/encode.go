package encoder

import (
	"io"
)

type Encoder interface {
	Encode(b io.Writer, v interface{}) error
}

type Decoder interface {
	Decode(b []byte) (interface{}, error)
}

type allocator func() interface{}
