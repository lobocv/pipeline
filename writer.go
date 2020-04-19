package pipeline

type PipeWriter interface {
	Write(payload []byte) error
}
