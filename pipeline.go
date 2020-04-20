package pipeline

import (
	"context"
	"fmt"
	"github.com/lobocv/pipeline/pencode"
	"sync"
)

type Processor interface {
	Process(ctx context.Context, payload interface{}) (interface{}, error)
}

func defaultErrorHandler(ctx context.Context, err error) error {
	fmt.Printf("An error was encountered in the pipeline: %s\n", err)
	return err
}

type errorHandler func(context.Context, error) error

type Tee struct {
}

type Pipeline struct {
	proc Processor

	// readers is a list of readers that will be read from for the input to the pipeline
	readers    []PipeReader
	readerLock sync.Mutex

	// writers is a list of writers that will be written to with the result payload at the end of the pipeline
	// TODO: Make writers map to encoders so that each writer can have it's own custom encoding but share encoded payloads.
	writers []PipeWriter

	// error handling function for pipeline errors
	errHandler func(context.Context, error) error

	// Encoder and decoder for the inputs and outputs of the pipeline
	enc pencode.Encoder
	dec pencode.Decoder

	// Done channel used to stop the pipeline if a fatal error occurs
	done chan struct{}
}

// NewPipeline creates a new pipeline
func NewPipeline(dec pencode.Decoder, enc pencode.Encoder) *Pipeline {
	return &Pipeline{enc: enc, dec: dec, errHandler: defaultErrorHandler, done: make(chan struct{})}
}

// AddReaders appends a reader to the input of this pipeline
func (p *Pipeline) AddReaders(r ...PipeReader) {
	p.readerLock.Lock()
	defer p.readerLock.Unlock()
	p.readers = append(p.readers, r...)
}

// AddWriters appends a writer to the output of this pipeline
func (p *Pipeline) AddWriters(w ...PipeWriter) {
	p.writers = append(p.writers, w...)
}

// SetProcessor sets the processor on the pipeline
func (p *Pipeline) SetProcessor(proc Processor) {
	p.proc = proc
}

// SetErrorHandler sets the processor on the pipeline
func (p *Pipeline) SetErrorHandler(h errorHandler) {
	p.errHandler = h
}

// Run is a blocking call that engages the pipeline
func (p *Pipeline) Run(ctx context.Context) {
	fmt.Println("Starting pipeline")
	for _, r := range p.readers {
		go p.listen(ctx, r)
	}
loop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Pipeline canceled")
			break loop
		case <-p.done:
			if len(p.readers) == 0 {
				fmt.Println("No more readers. Stopping pipeline")
				break loop
			}
		}
	}
	fmt.Println("Exiting pipeline gracefully")
}

func (p *Pipeline) RemoveReader(r PipeReader) {
	fmt.Println("Removing reader", r)
	n := 0
	for _, otherReader := range p.readers {
		if r != otherReader {
			p.readers[n] = otherReader
			n++
		}
	}
	p.readerLock.Lock()
	p.readers = p.readers[:n]
	p.readerLock.Unlock()
	fmt.Println("Readers remaining", p.readers)
}

// listen starts the pipelines for the specified reader
func (p *Pipeline) listen(ctx context.Context, r PipeReader) {
	var (
		errChan = make(chan error, len(p.readers))
	)
	fmt.Println("Starting reader")
loop:
	for {
		select {
		default:
			// Perform a blocking read on the PipeReader
			raw, err := r.Read()
			if err != nil {
				if err == EOF {
					fmt.Println("Reader reached EOF")
					p.RemoveReader(r)
					p.done <- struct{}{}
					break loop
				}
				errChan <- err
				continue
			}

			// Decode the byte stream into the payload
			dataPayload, err := p.dec.Decode(raw)
			if err != nil {
				errChan <- err
				continue
			}

			// Pass the payload to be processed
			result, err := p.proc.Process(ctx, dataPayload)
			if err != nil {
				errChan <- err
				continue
			}

			var rawResult []byte
			if rawResult, err = p.enc.Encode(result); err != nil {
				errChan <- err
				continue
			}

			// write the results of the payload
			if err = p.write(rawResult); err != nil {
				errChan <- err
				continue
			}

		case err := <-errChan:
			err = p.errHandler(ctx, err)
			if _, ok := err.(FatalError); ok {
				p.done <- struct{}{}
			}
		case <-ctx.Done():
			fmt.Println("Stopping reading from reader")
			break loop
		}
	}

	fmt.Println("Stopping reader")
}

func (p *Pipeline) AddEncoder(enc pencode.Encoder, _ PipeWriter) {
	// TODO: Make writers map to encoders so that each writer can have it's own custom encoding but share encoded payloads.
	p.enc = enc
}

// write implements PipeWriter as a multi-writer. It encodes and then writes the payload to all registered PipeWriters
func (p *Pipeline) write(results []byte) error {
	var errors []pipelineError

	for _, w := range p.writers {
		if _, err := w.Write(results); err != nil {
			// Create a pipelineError from the returned error
			var pipelineErr pipelineError
			pipelineErr.FromError(err)
			errors = append(errors, pipelineErr)
		}
	}
	if len(errors) > 0 {
		return overallError(errors...)
	}
	return nil
}
