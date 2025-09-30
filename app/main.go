package main

import (
	"fmt"
	"io"
	"os"

	"github.com/mmarchesotti/build-your-own-grep/app/nfasimulator"
)

// Usage: echo <input_text> | your_program.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a singular line
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		fmt.Println("Program returned 2")
		os.Exit(2)
	}

	ok := nfasimulator.Simulate(line, pattern)

	if !ok {
		fmt.Println("Program returned 1")
		os.Exit(1)
	}

	fmt.Println("Program returned 0")
	// default exit code is 0 which means success
}
