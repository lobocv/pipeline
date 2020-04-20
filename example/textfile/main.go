package main

import (
	"context"
	"fmt"
	"github.com/lobocv/pipeline/pencode"
	"github.com/lobocv/pipeline/pipeio"
	"math/rand"
	"strings"
	"time"

	"github.com/lobocv/pipeline"
	"os"
)

type ExampleLineProcessor struct{}

// Decorate the line with XXX. Sleep for a random amount of time to show asynchronicity
func (p ExampleLineProcessor) Process(ctx context.Context, payload interface{}) (interface{}, error) {
	s := string(payload.([]byte))
	sleep := rand.Intn(1000)
	time.Sleep(time.Duration(sleep) * time.Millisecond)
	processedLine := fmt.Sprintf("XXXX %s XXXX\n", strings.Trim(s, "\n"))
	return []byte(processedLine), nil
}

type ExampleWriter struct{}

func (p ExampleWriter) Write(payload []byte) (int, error) {
	fmt.Println("Writing:", string(payload))
	return len(payload), nil
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

	passthrough := pencode.PassThrough{}
	p := pipeline.NewPipeline(passthrough, passthrough)
	p.SetProcessor(ExampleLineProcessor{})
	p.AddReaders(reader1, reader2)
	p.AddWriters(ExampleWriter{}, os.Stdout, outFile)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	p.Run(ctx)

}

func mustSucceed(err error) {
	if err != nil {
		panic(err)
	}
}
