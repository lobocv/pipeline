package pipeio

import (
	"bufio"
	"os"
)

// FileReader is a reader parses a file based on a delimiter
type FileReader struct {
	r     *bufio.Reader
	delim byte
}

// Read up to the next delimiter
func (f FileReader) Read() ([]byte, error) {
	return f.r.ReadBytes(f.delim)
}

// NewFileReader creates a new file reader
func NewFileReader(path string, delim byte) (*FileReader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(f)
	return &FileReader{r: reader, delim: delim}, nil
}
