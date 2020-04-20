package pencode

import (
	"fmt"
)

type PassThrough struct {
}

func (d PassThrough) Decode(b []byte) (interface{}, error) {
	return b, nil
}

func (d PassThrough) Encode(v interface{}) ([]byte, error) {
	raw, ok := v.([]byte)
	if !ok {
		return nil, fmt.Errorf("passthrough expects []byte value but instead got %T", v)
	}
	return raw, nil
}
