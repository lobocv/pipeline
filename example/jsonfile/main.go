package main

import (
	"context"
	"fmt"
	"github.com/lobocv/pipeline/encoder"
	"github.com/lobocv/pipeline/pipeio"
	"math/rand"
	"time"

	"github.com/lobocv/pipeline"
	"os"
)

type ExampleJSONProcessor struct{}

// Decorate the line with XXX. Sleep for a random amount of time to show asynchronicity
func (p ExampleJSONProcessor) Process(payload interface{}) (interface{}, error) {

	s := payload.(*Input)
	sleep := rand.Intn(1000)
	time.Sleep(time.Duration(sleep) * time.Millisecond)
	processedLine := fmt.Sprintf("Processor sees: First Name = %s Last Name = %s\n", s.First, s.Last)
	return []byte(processedLine), nil
}

type ExampleWriter struct{}

func (p ExampleWriter) Write(payload []byte) error {
	fmt.Println("Writing:", string(payload))
	return nil
}

type Input struct {
	First string `json:"first_name"`
	Last  string `json:"last_name"`
}

func main() {
	rand.Seed(time.Now().UnixNano())
	reader1, err := pipeio.NewFileLineReader("./input.txt")
	mustSucceed(err)

	outFile, err := os.OpenFile("./output.txt", os.O_WRONLY|os.O_CREATE, 0600)
	mustSucceed(err)
	fileWriter := pipeio.NewFileWriter(outFile)

	passthrough := encoder.PassThrough{}

	alloc := func() interface{} {
		return new(Input)
	}

	p := pipeline.NewPipeline(encoder.NewJSONDecoder(alloc, false), passthrough)
	p.SetProcessor(ExampleJSONProcessor{})
	p.AddReaders(reader1)
	p.AddWriters(ExampleWriter{}, fileWriter)
	p.Run(context.Background())
}

func mustSucceed(err error) {
	if err != nil {
		panic(err)
	}
}
