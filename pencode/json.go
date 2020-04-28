package pencode

import (
	"bytes"
	"encoding/json"
)

// JSONDecoder decodes the byte payload into the given struct via the allocator
type JSONDecoder struct {
	Strict   bool
	allocate allocator
}

// NewJSONDecoder creates a new json decoder with the given allocator
func NewJSONDecoder(alloc allocator, strict bool) *JSONDecoder {
	return &JSONDecoder{allocate: alloc, Strict: strict}
}

// Decode decodes the payload into the allocated struct
func (d JSONDecoder) Decode(b []byte) (interface{}, error) {
	var v = d.allocate()
	reader := bytes.NewBuffer(b)
	dec := json.NewDecoder(reader)
	if d.Strict {
		dec.DisallowUnknownFields()
	}
	err := dec.Decode(v)
	return v, err
}

// JSONEncoder encodes structs into json payloads
type JSONEncoder struct{}

// NewJSONEncoder creates a new JSONEncoder
func NewJSONEncoder() *JSONEncoder {
	return &JSONEncoder{}
}

// Encode encodes the struct into a []byte payload
func (d JSONEncoder) Encode(v interface{}) ([]byte, error) {
	writer := bytes.Buffer{}
	enc := json.NewEncoder(&writer)
	err := enc.Encode(v)
	return writer.Bytes(), err
}
