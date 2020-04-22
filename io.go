package pipeline

import (
	"github.com/lobocv/pipeline/pencode"
	"io"
	"io/ioutil"
)

type pipeReader interface {
	Read() (interface{}, error)
}

type pipeWriter interface {
	Write(result interface{}) (int, error)
	Close() error
}

type MessageReader interface {
	Read() ([]byte, error)
}

type coupler struct {
	data      chan interface{}
	doneWrite chan struct{}
}

func newCoupler() *coupler {
	return &coupler{data: make(chan interface{}), doneWrite: make(chan struct{})}
}

func (c coupler) Write(result interface{}) (int, error) {
	c.data <- result
	return 0, nil
}

func (c coupler) Read() (interface{}, error) {
	select {
	case v := <-c.data:
		return v, nil
	case <-c.doneWrite:
		return nil, EOF
	}
}

func (c coupler) Close() error {
	c.doneWrite <- struct{}{}
	return nil
}

// bufferReader contains a io.Reader and a decoder and satisfies the pipeReader interface
type bufferReader struct {
	buf []byte
	r   io.Reader
	dec pencode.Decoder
}

func newBufferReader(r io.Reader, buf []byte, dec pencode.Decoder) *bufferReader {
	return &bufferReader{r: r, buf: buf, dec: dec}
}

func (b bufferReader) Read() (interface{}, error) {
	n, err := b.r.Read(b.buf)
	if err != nil {
		return nil, err
	}

	v, err := b.dec.Decode(b.buf[:n])
	if err != nil {
		return nil, err
	}

	return v, nil
}

// messageInput contains a MessageReader and a decoder and satisfies the pipeReader interface
type messageInput struct {
	r   MessageReader
	dec pencode.Decoder
}

func newMessageInput(r MessageReader, dec pencode.Decoder) *messageInput {
	return &messageInput{r: r, dec: dec}
}

// Read reads from the MessageReader and decodes the bytes
func (p *messageInput) Read() (interface{}, error) {
	// read the raw input
	raw, err := p.r.Read()
	if err != nil {
		return nil, err
	}
	// decode the input
	v, err := p.dec.Decode(raw)
	if err != nil {
		return nil, err
	}
	return v, nil
}

type pipeOutput struct {
	w   io.WriteCloser
	enc pencode.Encoder
}

// Write encodes the result and writes it to the io.Writer
func (p *pipeOutput) Write(result interface{}) (int, error) {
	var n int
	// encode the results
	raw, err := p.enc.Encode(result)
	if err != nil {
		return 0, err
	}

	// write the results of the payload
	if n, err = p.w.Write(raw); err != nil {
		return n, err
	}
	return n, nil
}

func (p *pipeOutput) Close() error {
	return p.w.Close()
}

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }

// NopWriteCloser returns a ReadCloser with a no-op Close method wrapping
// the provided Reader r.
func NopWriteCloser(wc io.Writer) io.WriteCloser {
	return nopCloser{Writer: wc}
}

func NopReadCloser(r io.Reader) io.ReadCloser {
	return ioutil.NopCloser(r)
}
