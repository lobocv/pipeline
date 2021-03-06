package generic

import (
	"io"
	"strings"
)

// EOF represents the end of file from a input stream. It is an alias for io.EOF
var EOF = io.EOF

type pipelineError struct {
	fatal     bool
	temporary bool
	msg       string
}

func (e *pipelineError) FromError(err error) {
	e.msg = err.Error()
	if specificErr, ok := err.(Fatal); ok {
		e.fatal = specificErr.Fatal()
	}
	if specificErr, ok := err.(TemporaryError); ok {
		e.fatal = specificErr.Temporary()
	}
}

func (e *pipelineError) Error() string {
	return e.msg
}

func overallError(errs ...error) *pipelineError {
	var overall pipelineError
	var msg strings.Builder
	_, _ = msg.Write([]byte("errors detected in the pipeline: ["))
	for ii, err := range errs {
		var pipelineErr pipelineError
		pipelineErr.FromError(err)
		// Combine error message strings
		_, _ = msg.Write([]byte(pipelineErr.msg))
		if ii < len(errs)-1 {
			_, _ = msg.Write([]byte("|"))
		}
		// Determine overall status of flags
		// If any errors are fatal, the overall error is fatal
		overall.fatal = overall.fatal || pipelineErr.fatal
		overall.temporary = overall.temporary && pipelineErr.temporary
	}
	_, _ = msg.Write([]byte("]"))
	overall.msg = msg.String()

	return &overall
}

// Fatal is an interface describing an error that is fatal to the pipeline
type Fatal interface {
	Fatal() bool
}

// FatalError is a basic implementation of a fatal error
type FatalError struct {
	error
}

// NewFatalError creates a new fatal error
func NewFatalError(error error) FatalError {
	return FatalError{error: error}
}

// Fatal indicates that this error is fatal
func (e FatalError) Fatal() bool {
	return true
}

// Temporary is an interface describing an error that is temporary and hence retry-able by the pipeline
type Temporary interface {
	Temporary() bool
}

// TemporaryError is a basic implementation of a temporary error
type TemporaryError struct {
	error
}

// NewTemporaryError creates a new temporary error
func NewTemporaryError(error error) TemporaryError {
	return TemporaryError{error: error}
}

// Temporary indicates that this error is temporary
func (e TemporaryError) Temporary() bool {
	return true
}
