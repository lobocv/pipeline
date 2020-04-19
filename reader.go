package pipeline

type PipeReader interface {
	Read() ([]byte, error)
}
