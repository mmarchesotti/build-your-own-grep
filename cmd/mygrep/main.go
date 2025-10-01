package main

import (
	"fmt"
	"io"
	"os"

	"github.com/mmarchesotti/build-your-own-grep/internal/nfasimulator"
)

// Usage: echo <input_text> | your_program.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	var line []byte
	if len(os.Args) == 3 {
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
			os.Exit(2)
		}
		line = content
	} else {
		content, err := os.ReadFile(os.Args[3])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: read file: %v\n", err)
			os.Exit(2)
		}
		line = content
	}

	ok, err := nfasimulator.Simulate(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: invalid pattern: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}
}
