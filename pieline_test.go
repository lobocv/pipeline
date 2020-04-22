package pipeline

import (
	"bytes"
	"context"
	"fmt"
	"github.com/lobocv/pipeline/pencode"
	. "github.com/smartystreets/goconvey/convey"
	"testing"

)

type Process struct {
	payloads [][]byte
}



func (p *Process) Process(ctx context.Context, payload interface{}) (interface{}, error) {
	raw := payload.([]byte)
	var payloadCopy = make([]byte, len(raw))
	copy(payloadCopy, raw)
	p.payloads = append(p.payloads, payloadCopy)
	return payload, nil
}

func TestPipeline(t *testing.T) {

	Convey("Given a pipeline", t, func() {
		p := NewPipeline()
		process := new(Process)
		p.SetProcessor(process)
		testInput := "This is test input"

		Convey(fmt.Sprintf("That reads from a filled buffer with text: '%s'", testInput), func() {
			input := bytes.NewBufferString(testInput)
			output := bytes.Buffer{}

			bufferSizes := []int{1, 2, len(testInput)/2, len(testInput)}

			for _, bufferSize := range bufferSizes {

				Convey(fmt.Sprintf("Where the reader has a buffer size of %d", bufferSize), func() {
					p.AddReader(input, pencode.PassThrough{}, make([]byte, bufferSize))

					Convey(fmt.Sprintf("And the pipeline writes to a different buffer"), func() {
						p.AddWriter(NopWriteCloser(&output), pencode.PassThrough{})

						nCalls := len(testInput) / bufferSize
						Convey(fmt.Sprintf("After running the pipeline, the output buffer should be filled and have made %d calls", nCalls), func() {
							p.Run(context.Background())
							So(output.String(), ShouldEqual, testInput)

							// Determine the number of calls to the processor that are made based on the reader buffer size
							var chunks [][]byte
							var chunk []byte
							for _, b := range []byte(testInput) {
								chunk = append(chunk, b)
								if len(chunk) == bufferSize {
									chunks = append(chunks, chunk)
									chunk = []byte{}
								}

							}
							So(process.payloads, ShouldResemble, chunks)
							So(process.payloads, ShouldHaveLength, nCalls)
						})
					})

				})
			}
		})

	})
}
