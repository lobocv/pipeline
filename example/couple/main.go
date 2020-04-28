package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/lobocv/pipeline"
	"github.com/lobocv/pipeline/pencode"
	"github.com/lobocv/pipeline/pipeio"
)

type person struct {
	Name string
}

type nameReadProcessor struct{}

// Read in the string and interpret it as a name for a person. Return a Person{}
func (p nameReadProcessor) Process(ctx context.Context, payload interface{}) (interface{}, error) {
	s := string(payload.([]byte))

	return person{Name: s}, nil
}

type personGreeter struct {
	greeting string
}

// Print a greeting to the supplied Person
func (p personGreeter) Process(ctx context.Context, payload interface{}) (interface{}, error) {
	person := payload.(person)

	return fmt.Sprintf("%s %s", p.greeting, person.Name), nil
}

func main() {
	var pipelines []*generic.Pipeline

	rand.Seed(time.Now().UnixNano())
	reader, err := pipeio.NewFileReader("./input.txt", '\n')
	mustSucceed(err)

	// Don't do any encoding / decoding
	passthrough := pencode.PassThrough{}

	p1 := generic.NewPipeline()
	pipelines = append(pipelines, p1)
	p1.AddMessageSource(reader, passthrough)

	// Set the processor
	p1.SetProcessor(nameReadProcessor{})

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for ii, greeting := range []string{"Hello there,", "Goodbye,"} {

		outFile, err := os.Create(fmt.Sprintf("./output%d.txt", ii+1))
		mustSucceed(err)

		p := generic.NewPipeline()
		pipelines = append(pipelines, p)

		p.SetProcessor(personGreeter{greeting: greeting})

		// Add stdout and a file as outputs
		p.AddWriter(generic.NopWriteCloser(os.Stdout), pencode.Printer{})
		p.AddWriter(outFile, pencode.Printer{})

		p1.Join(p)
	}
	generic.Run(ctx, pipelines...)
}

func mustSucceed(err error) {
	if err != nil {
		panic(err)
	}
}
