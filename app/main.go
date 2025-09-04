package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

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

func negativeMatchIndex(line []byte, pattern string) int {
	for index, character := range line {
		if !strings.ContainsRune(pattern, rune(character)) {
			return index
		}
	}
	return -1
}

var specialPatterns []string = []string{`\d`, `\w`}

func getSpecialPatternCharacterSet(specialPattern string) string {
	switch specialPattern {
	case `\d`:
		return "1234567890"
	case `\w`:
		return "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_"
	}
	return ""
}

func getStartingSpecialCharacter(pattern string) (bool, string) {
	for _, specialPattern := range specialPatterns {
		if strings.HasPrefix(pattern, specialPattern) {
			return true, specialPattern
		}
	}
	return false, ""
}

func matchLine(line []byte, pattern string) (bool, error) {
	if utf8.RuneCountInString(pattern) == 0 {
		return true, nil
	}

	if len(line) == 0 {
		return false, nil
	}

	var singleCharacterPattern string
	startsWithSpecialCharacter, startingSpecialCharacter := getStartingSpecialCharacter(pattern)

	if startsWithSpecialCharacter {
		singleCharacterPattern = startingSpecialCharacter
	} else if strings.HasPrefix(pattern, "[") {
		groupEndIndex := strings.Index(pattern, "]")
		if groupEndIndex < 0 {
			return false, fmt.Errorf("unsupported pattern: %q", pattern)
		}
		singleCharacterPattern = pattern[0 : groupEndIndex+1]
	} else {
		singleCharacterPattern = pattern[0:1]
	}

	matchIndex, _ := matchSingleCharacter(line, singleCharacterPattern)
	if matchIndex >= 0 {
		return matchLine(line[matchIndex+1:], pattern[len(singleCharacterPattern):])
	} else {
		return false, nil
	}
}

func matchSingleCharacter(line []byte, pattern string) (int, error) {

	var matchIndex int
	var characterSet string

	isSpecialPattern := contains(specialPatterns, pattern)
	isCharacterGroup := strings.HasPrefix(pattern, "[") && strings.HasSuffix(pattern, "]")
	isNegativeCharacterGroup := strings.HasPrefix(pattern, "[^") && strings.HasSuffix(pattern, "]")

	if isNegativeCharacterGroup {
		characterSet = pattern[2 : len(pattern)-1]
		matchIndex = negativeMatchIndex(line, characterSet)
	} else {
		if isSpecialPattern {
			characterSet = getSpecialPatternCharacterSet(pattern)
		} else if isCharacterGroup {
			characterSet = pattern[1 : len(pattern)-1]
		} else {
			characterSet = pattern
		}
		matchIndex = bytes.IndexAny(line, characterSet)
	}

	return matchIndex, nil
}
