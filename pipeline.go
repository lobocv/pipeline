package generic

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/lobocv/pipeline/pencode"
)

// Processor is an interface describing the processing component of the pipeline
type Processor interface {
	Process(ctx context.Context, payload interface{}) (interface{}, error)
}

type defaultErrorHandler struct {
}

func (h *defaultErrorHandler) HandleError(ctx context.Context, err error) error {
	fmt.Printf("An error was encountered in the pipeline: %s\n", err)
	return err
}

type errorHandler interface {
	HandleError(context.Context, error) error
}

// Logger is an interface for logging used in the pipeline
type Logger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	Error(format string, err error, v ...interface{})
}

type defaultLogger struct {
	log.Logger
}

func (l *defaultLogger) Error(format string, err error, v ...interface{}) {
	args := []interface{}{err}
	args = append(args, v...)
	l.Logger.Printf(format, args...)
}

// Pipeline represents a processing pattern where inputs are read, processed and written to outputs in a
// flexible and extensible manner. There can be multiple inputs and outputs that run concurrently at a time.
type Pipeline struct {
	log  Logger
	proc Processor

	// readers is a list of readers that will be read from for the input to the pipeline
	readers    []pipeReader
	readerLock sync.Mutex

	// writers is a list of writers that will be written to with the result payload at the end of the pipeline
	writers []pipeWriter

	// error handling function for pipeline errors
	errHandler errorHandler

	// Done channel used to stop the pipeline if a fatal error occurs
	done chan struct{}
}

// NewPipeline creates a new pipeline
func NewPipeline() *Pipeline {
	return &Pipeline{errHandler: &defaultErrorHandler{}, log: &defaultLogger{Logger: *log.New(os.Stderr, "", log.LstdFlags)}, done: make(chan struct{})}
}

// AddMessageSource appends a MessageReader to the input of this pipeline
func (p *Pipeline) AddMessageSource(r MessageReader, dec pencode.Decoder) {
	p.readerLock.Lock()
	defer p.readerLock.Unlock()
	p.readers = append(p.readers, newMessageInput(r, dec))
}

// AddReader appends an io.Reader to the input of this pipeline
func (p *Pipeline) AddReader(r io.Reader, dec pencode.Decoder, buf []byte) {
	p.readerLock.Lock()
	defer p.readerLock.Unlock()
	p.readers = append(p.readers, newBufferReader(r, buf, dec))
}

// RemoveReader removes the reader from the pipeline
func (p *Pipeline) RemoveReader(r pipeReader) {
	p.log.Println("Removing reader")
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
	p.log.Println("Readers remaining", p.readers)
}

// AddWriter appends a io.Writer to the output of this pipeline
func (p *Pipeline) AddWriter(w io.WriteCloser, enc pencode.Encoder) {
	p.writers = append(p.writers, &pipeOutput{w: w, enc: enc})
}

// Join joins the output of this pipeline to the input of the provided pipeline
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

// SetLogger sets the logger on the pipeline
func (p *Pipeline) SetLogger(l Logger) {
	p.log = l
}

// Run is a blocking call that engages the pipeline
func (p *Pipeline) Run(ctx context.Context) {
	p.log.Println("Starting pipeline")
	for _, r := range p.readers {
		go p.listen(ctx, r)
	}
loop:
	for {
		select {
		case <-ctx.Done():
			p.log.Println("Pipeline canceled")
			break loop
		case <-p.done:
			if len(p.readers) == 0 {
				p.log.Println("No more readers. Stopping pipeline")
				break loop
			}
		}
	}
	for _, w := range p.writers {
		p.log.Println("Closing writer")
		err := w.Close()
		if err != nil {
			p.log.Println("Error closing writer: ", err)
		}
		p.log.Println("Finished closing writer")

	}
	p.log.Println("Exiting pipeline gracefully")
}

// listen starts processing data for a given pipeReader
func (p *Pipeline) listen(ctx context.Context, r pipeReader) {
	var (
		errChan = make(chan error, len(p.readers))
	)
	p.log.Println("Starting reader")
loop:
	for {
		select {
		default:
			// Perform a blocking read on the pipeReader
			dataPayload, err := r.Read()
			if err != nil {
				if err == EOF {
					p.log.Println("Reader reached EOF")
					p.RemoveReader(r)
					p.done <- struct{}{}
					break loop
				}
				p.log.Error("Error during read: %s", err)
				errChan <- err
				continue
			}

			// Pass the payload to be processed
			result, err := p.proc.Process(ctx, dataPayload)
			if err != nil {
				p.log.Error("Error during processing: %s", err)
				errChan <- err
				continue
			}

			// write the results of the payload
			if err = p.write(result); err != nil {
				errChan <- err
				continue
			}

		case err := <-errChan:
			err = p.errHandler.HandleError(ctx, err)
			if _, ok := err.(FatalError); ok {
				p.done <- struct{}{}
			}
		case <-ctx.Done():
			p.log.Println("Stopping reading from reader")
			break loop
		}
	}

	p.log.Println("Stopping reader")
}

// write implements pipeWriter as a multi-writer. It encodes and then writes the payload to all registered PipeWriters
// This differs from io.MultiWriter because it does not stop writing on errors and instead returns a combined error
// for any failing writes.
func (p *Pipeline) write(results interface{}) error {
	var errors []error
	for _, w := range p.writers {
		if _, err := w.Write(results); err != nil {
			p.log.Error("Error during write: %s", err)
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return overallError(errors...)
	}
	return nil
}

// Run engages all the provided pipelines. This is useful for when multiple pipelines are coupled together
// This function blocks until all pipelines have finished running
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
