package main

import (
	"flag"
	"fmt"

	"github.com/sceptero/house-extractor/internal/extractor"
)

func main() {
	var inputFilePath, outputFilePath string
	flag.StringVar(&inputFilePath, "i", "", "an input file path")
	flag.StringVar(&outputFilePath, "o", "./output.lua", "an output file path")
	flag.Parse()

	if inputFilePath == "" {
		fmt.Println("Please specify input file path (use -i flag)")
		return
	}

	fmt.Printf("Extracting house data from file `%s` into file `%s`\n", inputFilePath, outputFilePath)

	e, err := extractor.New(inputFilePath, outputFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = e.Do()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Success")
}
