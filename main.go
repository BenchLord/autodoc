package main

import (
	"autodoc/betterfile"
	"fmt"
	"io"
)

func main() {

	f, err := betterfile.NewFile("test.proto")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	lines := make([][]byte, 0)
	for {
		line, err := f.GetLineAsBytes()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		lines = append(lines, line)
	}

	for _, value := range lines {
		fmt.Print(string(value))
	}
}
