package main

import (
	"autodoc/protofile"
	"fmt"
)

func main() {

	_, err := protofile.NewFile("test.proto")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
