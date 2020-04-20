package pipeline

import "io"

type PipeWriter interface {
	io.Writer
}
