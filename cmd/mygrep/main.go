package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/mmarchesotti/build-your-own-grep/internal/nfasimulator"
)

const usage = `Usage: mygrep [options] <pattern> [path...]

Search for PATTERN in each PATH. If no PATH is provided,
the search reads from standard input.

Options:
  -r    Recursively search subdirectories. When this flag is used,
        the trailing path must be a single directory.

Examples:
  mygrep 'apple' file1.txt file2.txt
  cat file.txt | mygrep 'apple'
  mygrep -r 'apple' ./my_project`

func main() {
	recursive := flag.Bool("r", false, "Recursive search")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "error: missing pattern")
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}

	pattern := args[0]
	paths := args[1:]

	matchFound := false
	var filenames []string
	if *recursive {
		if len(paths) != 1 {
			fmt.Fprintln(os.Stderr, "error: recursive search requires exactly one directory path")
			os.Exit(2)
		}
		directory := paths[0]

		err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Printf("error accessing path %q: %v\n", path, err)
				return err
			}

			if !d.IsDir() {
				filenames = append(filenames, path)
			}

			return nil
		})

		if err != nil {
			fmt.Printf("error walking the path %q: %v\n", directory, err)
		}

	} else {
		filenames = paths
	}

	if len(filenames) == 0 {
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
		for _, filename := range filenames {
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
