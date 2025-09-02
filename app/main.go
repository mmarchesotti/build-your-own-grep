package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

// Ensures gofmt doesn't remove the "bytes" import above (feel free to remove this!)
var _ = bytes.ContainsAny

// Usage: echo <input_text> | your_program.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok, err := matchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

func contains(slice []string, target string) bool {
	for _, element := range slice {
		if element == target {
			return true
		}
	}
	return false
}

func isNotContained(line []byte, pattern string) bool {
	for _, character := range line {
		if !strings.ContainsRune(pattern, rune(character)) {
			return true
		}
	}
	return false
}

func matchLine(line []byte, pattern string) (bool, error) {
	var ok bool
	specialPatterns := []string{`\d`, `\w`}

	isSingleCharacter := utf8.RuneCountInString(pattern) == 1
	isSpecialPattern := contains(specialPatterns, pattern)
	isCharacterGroup := strings.HasPrefix(pattern, "[") && strings.HasSuffix(pattern, "]")
	isNegativeCharacterGroup := strings.HasPrefix(pattern, "[^") && strings.HasSuffix(pattern, "]")

	if !(isSingleCharacter || isSpecialPattern || isCharacterGroup) {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}

	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if isNegativeCharacterGroup {
		pattern = pattern[2 : len(pattern)-1]
		ok = isNotContained(line, pattern)
	} else {
		if pattern == `\d` {
			pattern = "1234567890"
		} else if pattern == `\w` {
			pattern = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_"
		} else if isCharacterGroup {
			pattern = pattern[1 : len(pattern)-1]
		}
		ok = bytes.ContainsAny(line, pattern)
	}

	return ok, nil
}
