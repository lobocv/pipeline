package generic

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/lobocv/pipeline/mocks"
)

var pCount int

// P is a test payload that goes through the pipeline
type P struct {
	id        int
	processed bool
}

// NewP creates a new P with an monotonically increasing ID for uniqueness
func NewP() P {
	pCount++
	return P{id: pCount}
}

// Payload is a grouping of the raw and processed form of the payload P
type Payload struct {
	raw  P
	proc P
}

type PipelineSuite struct {
	suite.Suite
	pipeline *Pipeline

	mockProc       *mocks.Processor
	mockWriters    []*mocks.PipeWriter
	mockReaders    []*mocks.PipeReader
	mockErrHandler *mocks.ErrorHandler
}

// Setup the test suite with a pipeline configured with a mock processor, error handler and one mock reader and writer
func (t *PipelineSuite) SetupTest() {
	t.mockProc = &mocks.Processor{}
	t.mockErrHandler = &mocks.ErrorHandler{}

	t.pipeline = NewPipeline()
	t.pipeline.SetProcessor(t.mockProc)
	t.pipeline.SetErrorHandler(t.mockErrHandler)
	t.addMockReader()
	t.addMockWriter()
}

func (t *PipelineSuite) TearDownTest() {
	t.mockErrHandler.AssertExpectations(t.T())
	t.mockProc.AssertExpectations(t.T())
	for _, m := range t.mockReaders {
		m.AssertExpectations(t.T())

	}
	for _, m := range t.mockWriters {
		m.AssertExpectations(t.T())
	}
	// Reset list of readers and writers
	t.mockWriters = nil
	t.mockReaders = nil

}

func (t *PipelineSuite) addMockReader() *mocks.PipeReader {
	m := &mocks.PipeReader{}
	t.mockReaders = append(t.mockReaders, m)
	t.pipeline.readers = append(t.pipeline.readers, m)
	return m
}

func (t *PipelineSuite) addMockWriter() *mocks.PipeWriter {
	m := &mocks.PipeWriter{}
	t.mockWriters = append(t.mockWriters, m)
	t.pipeline.writers = append(t.pipeline.writers, m)
	return m
}

func generatePayloads(n int) (payloads []Payload) {
	for ii := 0; ii < n; ii++ {
		raw := NewP()
		proc := raw
		proc.processed = true
		payloads = append(payloads, Payload{raw: raw, proc: proc})
	}
	return
}

// This test checks that the pipeline works for a single reader and writer
func (t *PipelineSuite) TestSingleReaderWriter() {
	mockReader := t.mockReaders[0]
	mockWriter := t.mockWriters[0]
	ctx := context.Background()
	numPayloads := 3
	payloads := []P{}
	for ii := 0; ii < numPayloads; ii++ {
		p := NewP()
		payloads = append(payloads, p)
		mockReader.On("Read").Return(p, nil).Once()
		result := p
		result.processed = true
		t.mockProc.On("Process", ctx, p).Return(result, nil).Once()
		mockWriter.On("Write", result).Return(0, nil).Once()
	}
	t.setMockEOF()

	t.pipeline.Run(ctx)
}

// This test checks that the pipeline works as expected when there are multiple readers and writers
func (t *PipelineSuite) TestMultiReaderWriter() {
	// Add additional readers and writers
	for ii := 0; ii < 2; ii++ {
		t.addMockReader()
		t.addMockWriter()
	}
	t.Len(t.pipeline.readers, 3)
	t.Len(t.pipeline.writers, 3)

	ctx := context.Background()
	for _, payload := range generatePayloads(3) {

		for _, m := range t.mockReaders {
			m.On("Read").Return(payload.raw, nil).Once()

			t.mockProc.On("Process", ctx, payload.raw).Return(payload.proc, nil).Once()

			for _, m := range t.mockWriters {
				m.On("Write", payload.proc).Return(0, nil).Once()
			}
		}
	}
	t.setMockEOF()
	t.pipeline.Run(ctx)
}

// This test checks that the pipeline works as expected when there are multiple readers and writers
func (t *PipelineSuite) TestPipelineErrorHandling() {
	t.addMockWriter()

	ctx := context.Background()
	mockReader := t.mockReaders[0]

	type writeError struct {
		err    error
		writer *mocks.PipeWriter
	}
	testCases := []struct {
		readErr, procErr error
		writeErr         []writeError
	}{
		{readErr: fmt.Errorf("read error from test")},
		{procErr: fmt.Errorf("proc error from test")},
		{writeErr: []writeError{
			{err: fmt.Errorf("write error from test"), writer: t.mockWriters[1]},
		}},
		{writeErr: []writeError{
			{err: fmt.Errorf("write error 1 from test"), writer: t.mockWriters[0]},
			{err: fmt.Errorf("write error 2 from test"), writer: t.mockWriters[1]},
		}},
	}

	for _, tc := range testCases {
		payload := generatePayloads(1)[0]
		if tc.readErr != nil {
			mockReader.On("Read").Return(nil, tc.readErr).Once()
			t.mockErrHandler.On("HandleError", ctx, tc.readErr).Return(tc.readErr)
			continue
		} else {
			mockReader.On("Read").Return(payload.raw, nil).Once()
		}

		if tc.procErr != nil && tc.readErr == nil {
			t.mockProc.On("Process", ctx, payload.raw).Return(nil, tc.procErr).Once()
			t.mockErrHandler.On("HandleError", ctx, tc.procErr).Return(tc.procErr)
			continue
		} else {
			t.mockProc.On("Process", ctx, payload.raw).Return(payload.proc, nil).Once()
		}

		for _, mockWriter := range t.mockWriters {

			var found bool
			for _, writeErr := range tc.writeErr {
				if mockWriter == writeErr.writer {
					mockWriter.On("Write", payload.proc).Return(0, writeErr.err).Once()
					found = true
				}
			}
			if !found {
				mockWriter.On("Write", payload.proc).Return(0, nil).Once()
			}
		}

		// Set the expected combined write error
		if len(tc.writeErr) > 0 {
			var errors []error
			for _, writerErr := range tc.writeErr {
				errors = append(errors, writerErr.err)
			}
			err := overallError(errors...)
			t.mockErrHandler.On("HandleError", ctx, err).Return(err)
		}

	}

	t.setMockEOF()
	t.pipeline.Run(ctx)
}

// setMockEOF sets the readers to all return EOF on their next call and the writers to expect a call to Close()
func (t *PipelineSuite) setMockEOF() {
	for _, m := range t.mockReaders {
		m.On("Read").Return(nil, EOF)
	}
	for _, m := range t.mockWriters {
		m.On("Close").Return(nil)
	}
}

func TestPipeline(t *testing.T) {
	s := PipelineSuite{}
	suite.Run(t, &s)
}
