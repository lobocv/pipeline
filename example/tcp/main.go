package main

import (
	"context"
	"fmt"
	"github.com/lobocv/pipeline/pencode"
	"net"
	"os"
	"time"

	"github.com/lobocv/pipeline"
)

type ExampleJSONProcessor struct{}

// Simple processor that reverses the input. Expects []byte payload
func (p ExampleJSONProcessor) Process(ctx context.Context, payload interface{}) (interface{}, error) {
	s := payload.([]byte)
	for ii := 0; ii < len(s)/2; ii++ {
		jj := len(s) - 1 - ii
		left := s[ii]
		s[ii] = s[jj]
		s[jj] = left
	}

	return payload, nil
}

func main() {
	port := 5001

	// Listen for tcp connections
	tcpListen, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 5001))
	mustSucceed(err)
	fmt.Printf("Starting listening on port %d: Use `netcat localhost %d` and type in messages\n", port, port)
	conn, err := tcpListen.Accept()
	mustSucceed(err)
	defer conn.Close()

	// Don't do any encoding
	passthrough := pencode.PassThrough{}

	// Create a new pipeline
	p := pipeline.NewPipeline()

	// Set the processor
	p.SetProcessor(ExampleJSONProcessor{})

	// Add the connection as both a reader and a writer
	p.AddReader(conn, passthrough, make([]byte, 1000))
	p.AddWriter(conn, passthrough)
	// Add standard out so we can see messages on the pipeline side
	p.AddWriter(os.Stdout, passthrough)

	// Start the pipeline. This is blocking so we can set a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	p.Run(ctx)
}

func mustSucceed(err error) {
	if err != nil {
		panic(err)
	}
}
