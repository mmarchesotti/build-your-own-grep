package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/mmarchesotti/build-your-own-grep/internal/nfasimulator"
)

func main() {
	if len(os.Args) < 3 || len(os.Args) > 4 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern> [file]\n")
		os.Exit(2)
	}

	pattern := os.Args[2]
	var input io.Reader

	if len(os.Args) == 4 {
		file, err := os.Open(os.Args[3])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: could not read file %s: %v\n", os.Args[3], err)
			os.Exit(2)
		}
		defer file.Close()
		input = file
	} else {
		input = os.Stdin
	}

	matchFound, err := processLines(input, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !matchFound {
		os.Exit(1)
	}
}

func processLines(input io.Reader, pattern string) (bool, error) {
	scanner := bufio.NewScanner(input)
	anyMatchFound := false

	for scanner.Scan() {
		line := scanner.Bytes()
		lineCopy := make([]byte, len(line))
		copy(lineCopy, line)

		ok, err := nfasimulator.Simulate(lineCopy, pattern)
		if err != nil {
			return false, fmt.Errorf("invalid pattern: %w", err)
		}

		if ok {
			anyMatchFound = true
			var out bytes.Buffer
			out.Write(line)
			out.WriteByte('\n')
			os.Stdout.Write(out.Bytes())
		}
	}

	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("error reading input: %w", err)
	}

	return anyMatchFound, nil
}
