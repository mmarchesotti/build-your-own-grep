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
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern> [file...]\n")
		os.Exit(2)
	}

	pattern := os.Args[2]
	matchFound := false
	if len(os.Args) == 3 {
		hasMatch, matchedLines, err := processLines(os.Stdin, pattern)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(2)
		}
		matchFound = hasMatch
		for _, line := range matchedLines {
			var out bytes.Buffer
			out.Write(line)
			out.WriteByte('\n')
			os.Stdout.Write(out.Bytes())
		}
	} else {
		for _, filename := range os.Args[3:] {
			file, err := os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: could not read file %s: %v\n", filename, err)
				os.Exit(2)
			}
			defer file.Close()

			hasMatch, matchedLines, err := processLines(file, pattern)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(2)
			}
			matchFound = matchFound || hasMatch
			for _, line := range matchedLines {
				var out bytes.Buffer
				out.Write([]byte(filename))
				out.WriteByte(':')
				out.Write(line)
				out.WriteByte('\n')
				os.Stdout.Write(out.Bytes())
			}
		}
	}

	if !matchFound {
		os.Exit(1)
	}
}

func processLines(input io.Reader, pattern string) (bool, [][]byte, error) {
	scanner := bufio.NewScanner(input)
	anyMatchFound := false

	var matchedLines [][]byte
	for scanner.Scan() {
		line := scanner.Bytes()
		lineCopy := make([]byte, len(line))
		copy(lineCopy, line)

		ok, err := nfasimulator.Simulate(lineCopy, pattern)
		if err != nil {
			return false, nil, fmt.Errorf("invalid pattern: %w", err)
		}

		if ok {
			anyMatchFound = true
			matchedLines = append(matchedLines, lineCopy)
		}
	}

	if err := scanner.Err(); err != nil {
		return false, nil, fmt.Errorf("error reading input: %w", err)
	}

	return anyMatchFound, matchedLines, nil
}
