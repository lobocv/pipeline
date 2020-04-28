package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/lobocv/pipeline"
	"github.com/lobocv/pipeline/pencode"
	"github.com/lobocv/pipeline/pipeio"
)

type exampleLineProcessor struct{}

// Decorate the line with XXX. Sleep for a random amount of time to show asynchronicity
func (p exampleLineProcessor) Process(ctx context.Context, payload interface{}) (interface{}, error) {
	s := string(payload.([]byte))
	sleep := rand.Intn(1000)
	time.Sleep(time.Duration(sleep) * time.Millisecond)
	processedLine := fmt.Sprintf("XXXX %s XXXX\n", strings.Trim(s, "\n"))
	return []byte(processedLine), nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	reader1, err := pipeio.NewFileReader("./input.txt", '\n')
	mustSucceed(err)

	reader2, err := pipeio.NewFileReader("./input2.txt", '\n')
	mustSucceed(err)

	outFile, err := os.Create("./output.txt")
	mustSucceed(err)

	// Don't do any encoding / decoding
	passthrough := pencode.PassThrough{}

	p := generic.NewPipeline()

	// Set the processor
	p.SetProcessor(exampleLineProcessor{})

	// Add the file readers and message sources
	p.AddMessageSource(reader1, passthrough)
	p.AddMessageSource(reader2, passthrough)

	// Add stdout and a file as outputs
	p.AddWriter(os.Stdout, passthrough)
	p.AddWriter(outFile, passthrough)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	p.Run(ctx)

}

func mustSucceed(err error) {
	if err != nil {
		panic(err)
	}
}
