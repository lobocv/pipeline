package pipeline

import (
	"context"
	"fmt"
	"github.com/lobocv/pipeline/pencode"
	"io"
	"sync"
)

type Processor interface {
	Process(ctx context.Context, payload interface{}) (interface{}, error)
}

func defaultErrorHandler(ctx context.Context, err error) error {
	panic(err)
	fmt.Printf("An error was encountered in the pipeline: %s\n", err)

	return err
}

type errorHandler func(context.Context, error) error

type Pipeline struct {
	proc Processor

	// readers is a list of readers that will be read from for the input to the pipeline
	readers    []pipeReader
	readerLock sync.Mutex

	// writers is a list of writers that will be written to with the result payload at the end of the pipeline
	writers []pipeWriter

	// error handling function for pipeline errors
	errHandler func(context.Context, error) error

	// Done channel used to stop the pipeline if a fatal error occurs
	done chan struct{}
}

// NewPipeline creates a new pipeline
func NewPipeline() *Pipeline {
	return &Pipeline{errHandler: defaultErrorHandler, done: make(chan struct{})}
}

// AddMessageSource appends a MessageReader to the input of this pipeline
func (p *Pipeline) AddMessageSource(r MessageReader, dec pencode.Decoder) {
	p.readerLock.Lock()
	defer p.readerLock.Unlock()
	p.readers = append(p.readers, newMessageInput(r, dec))
}

// AddMessageSource appends an io.Reader to the input of this pipeline
func (p *Pipeline) AddReader(r io.Reader, dec pencode.Decoder, buf []byte) {
	p.readerLock.Lock()
	defer p.readerLock.Unlock()
	p.readers = append(p.readers, newBufferReader(r, buf, dec))
}

// AddWriter appends a io.Writer to the output of this pipeline
func (p *Pipeline) AddWriter(w io.WriteCloser, enc pencode.Encoder) {
	p.writers = append(p.writers, &pipeOutput{w: w, enc: enc})
}

func (p *Pipeline) Join(out *Pipeline) {
	p.readerLock.Lock()
	defer p.readerLock.Unlock()
	c := newCoupler()
	out.readers = append(out.readers, c)
	p.writers = append(p.writers, c)
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
	for _, w := range p.writers {
		fmt.Println("Closing writer", w)
		err := w.Close()
		if err != nil {
			fmt.Println("Error closing writer: ", err)
		}
		fmt.Println("Finished closing writer", w)

	}
	fmt.Println("Exiting pipeline gracefully")
}

func (p *Pipeline) RemoveReader(r pipeReader) {
	fmt.Println("Removing reader")
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
func (p *Pipeline) listen(ctx context.Context, r pipeReader) {
	var (
		errChan = make(chan error, len(p.readers))
	)
	fmt.Println("Starting reader")
loop:
	for {
		select {
		default:
			// Perform a blocking read on the pipeReader
			dataPayload, err := r.Read()
			if err != nil {
				if err == EOF {
					fmt.Println("Reader reached EOF")
					p.RemoveReader(r)
					p.done <- struct{}{}
					break loop
				}
				fmt.Println("ERROR on read: ", err)
				errChan <- err
				continue
			}

			// Pass the payload to be processed
			result, err := p.proc.Process(ctx, dataPayload)
			if err != nil {
				fmt.Println("ERROR on process: ", err)
				errChan <- err
				continue
			}

			// write the results of the payload
			if err = p.write(result); err != nil {
				fmt.Println("ERROR on write: ", err)
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

// write implements pipeWriter as a multi-writer. It encodes and then writes the payload to all registered PipeWriters
func (p *Pipeline) write(results interface{}) error {
	var errors []pipelineError
	for _, w := range p.writers {
		if _, err := w.Write(results); err != nil {
			fmt.Println("ERROR on write: ", err)
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

func Run(ctx context.Context, pipelines ...*Pipeline) {
	wg := sync.WaitGroup{}

	for _, p := range pipelines {
		wg.Add(1)
		go func(p *Pipeline) {
			defer wg.Done()
			p.Run(ctx)
		}(p)
	}
	wg.Wait()
}
