package main

import (
	"context"
	"fmt"
	"github.com/lobocv/pipeline/pencode"
	"github.com/lobocv/pipeline/pipeio"
	"math/rand"
	"time"

	"github.com/lobocv/pipeline"
	"os"
)

type Person struct {
	Name string
}

type NameReadProcessor struct{}

// Decorate the line with XXX. Sleep for a random amount of time to show asynchronicity
func (p NameReadProcessor) Process(ctx context.Context, payload interface{}) (interface{}, error) {
	s := string(payload.([]byte))

	return Person{Name: s}, nil
}

type PersonGreeter struct {
	greeting string
}

// Decorate the line with XXX. Sleep for a random amount of time to show asynchronicity
func (p PersonGreeter) Process(ctx context.Context, payload interface{}) (interface{}, error) {
	person := payload.(Person)

	return fmt.Sprintf("%s %s", p.greeting, person.Name), nil
}

func main() {
	var pipelines []*pipeline.Pipeline

	rand.Seed(time.Now().UnixNano())
	reader, err := pipeio.NewFileReader("./input.txt", '\n')
	mustSucceed(err)

	// Don't do any encoding / decoding
	passthrough := pencode.PassThrough{}

	p1 := pipeline.NewPipeline()
	pipelines = append(pipelines, p1)
	p1.AddMessageSource(reader, passthrough)

	// Set the processor
	p1.SetProcessor(NameReadProcessor{})

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for ii, greeting := range []string{"Hello there,", "Goodbye,"} {

		outFile, err := os.Create(fmt.Sprintf("./output%d.txt", ii+1))
		mustSucceed(err)

		p := pipeline.NewPipeline()
		pipelines = append(pipelines, p)

		p.SetProcessor(PersonGreeter{greeting: greeting})

		// Add stdout and a file as outputs
		p.AddWriter(pipeline.NopWriteCloser(os.Stdout), pencode.Printer{})
		p.AddWriter(outFile, pencode.Printer{})

		p1.Join(p)
	}
	pipeline.Run(ctx, pipelines...)
}

func mustSucceed(err error) {
	if err != nil {
		panic(err)
	}
}
