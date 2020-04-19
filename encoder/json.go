package encoder

import (
	"bytes"
	"encoding/json"
)

type JSONDecoder struct {
	Strict   bool
	allocate allocator
}

func NewJSONDecoder(alloc allocator, strict bool) *JSONDecoder {
	return &JSONDecoder{allocate: alloc, Strict: strict}
}

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
