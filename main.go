package main

import (
	"fmt"
	"log"
	"os"
	"path"
)

func main() {

	if len(os.Args) == 1 {
		log.Println("Incorrect usage, input file must be provided as the first argument")
		log.Fatalf("%s [input-file] [output-file]\n", os.Args[0]) // For debugging purposes
	} else if len(os.Args) == 2 {
		log.Println("Incorrect usage, output file must be provided as the second argument")
		log.Fatalf("%s [input-file] [output-file]\n", os.Args[0]) // For debugging purposes
	}

	filePath := path.Clean(os.Args[1])
	outputPath := path.Clean(os.Args[2])

	body, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Couldn't read file: %v", filePath)
	}

  fmt.Printf("Input file: %s\n", filePath)
	tokenizer := NewTokenizer(string(body))
	emitter := NewEmiiter()
	parser := NewParser(tokenizer, emitter)
	parser.Parse()

	emitter.SaveToFile(outputPath)
  fmt.Printf("Output file: %s\n", outputPath)
}
