package pipeio

import (
	"bufio"
	"os"
)

type FileLineReader struct {
	r *bufio.Reader
}

func (f FileLineReader) Read() ([]byte, error) {
	return f.r.ReadBytes('\n')
}

func NewFileLineReader(path string) (*FileLineReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	return &FileLineReader{r: reader}, nil
}
