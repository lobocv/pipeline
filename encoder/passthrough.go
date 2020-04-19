package encoder

import "io"

type PassThrough struct {
}

func (d PassThrough) Decode(b []byte) (interface{}, error) {
	return b, nil
}

func (d PassThrough) Encode(b io.Writer, v interface{}) error {
	_, err := b.Write(v.([]byte))
	return err
}
