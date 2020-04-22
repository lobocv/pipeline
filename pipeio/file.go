package pipeio

import (
	"bufio"
	"os"
)

type FileReader struct {
	r     *bufio.Reader
	delim byte
}

func (f FileReader) Read() ([]byte, error) {
	return f.r.ReadBytes(f.delim)
}

func NewFileReader(path string, delim byte) (*FileReader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(f)
	return &FileReader{r: reader, delim: delim}, nil
}
