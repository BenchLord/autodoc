package main

import (
	"autodoc/protofile"
	"encoding/json"
	"fmt"
)

func main() {

	pf, err := protofile.NewFile("test.proto")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	messages := pf.GetMessages()

	j, _ := json.MarshalIndent(messages, "", "  ")
	fmt.Printf("%s\n", j)
}
