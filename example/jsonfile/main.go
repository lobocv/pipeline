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

type exampleJSONProcessor struct{}

type output struct {
	input
	Title string `json:"title"`
}

// Simple processor that adds a title field to the input data struct
func (p exampleJSONProcessor) Process(ctx context.Context, payload interface{}) (interface{}, error) {
	sleep := rand.Intn(1000)
	time.Sleep(time.Duration(sleep) * time.Millisecond)

	s := payload.(*input)
	fmt.Printf("Processor sees: First Name = %s Last Name = %s. Adding Mr title.\n", s.First, s.Last)

	return output{input: *s, Title: "Mr"}, nil
}

type input struct {
	First string `json:"first_name"`
	Last  string `json:"last_name"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Read from the files one line at a time
	lineReader, err := pipeio.NewFileReader("./input.txt", '\n')
	mustSucceed(err)

	lineReader2, err := pipeio.NewFileReader("./input2.txt", '\n')
	mustSucceed(err)

	// Open a file to write output to
	outFile, err := os.Create("./output.txt")
	mustSucceed(err)

	// Define the allocator, JSON encoder and decoder
	alloc := func() interface{} {
		return new(input)
	}
	dec := pencode.NewJSONDecoder(alloc, false)
	enc := pencode.NewJSONEncoder()

	// Create a new pipeline
	p := generic.NewPipeline()

	// Set the processor
	p.SetProcessor(exampleJSONProcessor{})

	// Add reads and writers
	p.AddMessageSource(lineReader, dec)
	p.AddMessageSource(lineReader2, dec)

	// Set writer to the output file
	p.AddWriter(outFile, enc)

	// Stat the pipeline. This is blocking so we can set a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	p.Run(ctx)
}

func mustSucceed(err error) {
	if err != nil {
		panic(err)
	}
}
