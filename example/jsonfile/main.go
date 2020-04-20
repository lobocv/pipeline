package main

import (
	"context"
	"fmt"
	"github.com/lobocv/pipeline/pencode"
	"github.com/lobocv/pipeline/pipeio"
	"math/rand"
	"net"
	"time"

	"github.com/lobocv/pipeline"
	"os"
)

type ExampleJSONProcessor struct{}

type Output struct {
	Input
	Title string `json:"title"`
}

// Simple processor that adds a title field to the input data struct
func (p ExampleJSONProcessor) Process(ctx context.Context, payload interface{}) (interface{}, error) {
	sleep := rand.Intn(1000)
	time.Sleep(time.Duration(sleep) * time.Millisecond)

	s := payload.(*Input)
	fmt.Printf("Processor sees: First Name = %s Last Name = %s. Adding Mr title.\n", s.First, s.Last)

	return Output{Input: *s, Title: "Mr"}, nil
}

type Input struct {
	First string `json:"first_name"`
	Last  string `json:"last_name"`
}

func main() {
	rand.Seed(time.Now().UnixNano())
	reader1, err := pipeio.NewFileLineReader("./input.txt")
	mustSucceed(err)
	reader2, err := pipeio.NewFileLineReader("./input2.txt")
	mustSucceed(err)

	outFile, err := os.OpenFile("./output.txt", os.O_WRONLY|os.O_CREATE, 0600)
	mustSucceed(err)

	// Use netcat -l -k localhost 8999 to see results :)
	tcpConn, tcpConnErr := net.Dial("tcp", "localhost:8999")
	// Use netcat -l -u -k localhost 8998 to see results :)
	udpConn, udpConnErr := net.Dial("udp", "localhost:8998")

	// Define the allocator, encoder and decoder
	alloc := func() interface{} {
		return new(Input)
	}
	dec := pencode.NewJSONDecoder(alloc, false)
	enc := pencode.NewJSONEncoder()

	// Create a new pipeline
	p := pipeline.NewPipeline(dec, enc)
	// Set the processor
	p.SetProcessor(ExampleJSONProcessor{})

	// Add reads and writers
	p.AddReaders(reader1, reader2)
	p.AddWriters(outFile, os.Stdout)
	// Only add network writers if there are listeners, otherwise you get pipeline errors
	if tcpConnErr == nil {
		p.AddWriters(tcpConn)
	}
	if udpConnErr == nil {
		p.AddWriters(udpConn)
	}

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
