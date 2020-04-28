package pencode

import (
	"fmt"
)

// PassThrough is both an encoder and decoder which acts as a passthrough (does not alter the payload)
type PassThrough struct {
}

// Decode does nothing with the data going through
func (d PassThrough) Decode(b []byte) (interface{}, error) {
	return b, nil
}

// Encode does nothing with the data going through, although it assumes it to be []byte
func (d PassThrough) Encode(v interface{}) ([]byte, error) {
	raw, ok := v.([]byte)
	if !ok {
		return nil, fmt.Errorf("passthrough expects []byte value but instead got %T", v)
	}
	return raw, nil
}

// Printer is an encoder that uses fmt.Sprintf to encode to []byte
type Printer struct {
}

// Encode encodes the payload into []byte using fmt.Sprintf
func (p Printer) Encode(v interface{}) ([]byte, error) {
	return []byte(fmt.Sprint(v)), nil
}
