package main

import (
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

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a singular line
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		fmt.Println("Program returned 2")
		os.Exit(2)
	}

	ok := matchAtBeginning(line, pattern)

	if !ok {
		fmt.Println("Program returned 1")
		os.Exit(1)
	}

	fmt.Println("Program returned 0")
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

func matchAtBeginning(line []byte, pattern string) bool {	
	hasStartOfStringAnchor := strings.HasPrefix(pattern, "^")
	if (hasStartOfStringAnchor) {
		pattern = pattern[1:]
	}

	for true {
		if (matchAtCurrentPoint(line, pattern)) {
			return true
		}
		line = line[1:]
		if (len(line) == 0 || hasStartOfStringAnchor) {
			break
		}
	}
	return false
}

func matchAtCurrentPoint(line []byte, pattern string) bool {
	if utf8.RuneCountInString(pattern) == 0 {
		return true
	}
	if hasEndOfStringAchor {
		return len(line) == 0
	}
	if len(line) == 0 {
		return false
	}

	var singularPattern string
	startsWithSpecialCharacter, startingSpecialCharacter := getStartingSpecialCharacter(pattern)

	if startsWithSpecialCharacter {
		singularPattern = startingSpecialCharacter
	} else if strings.HasPrefix(pattern, "[") {
		groupEndIndex := strings.Index(pattern, "]")
		if groupEndIndex < 0 {
			return false
		}
		singularPattern = pattern[0 : groupEndIndex+1]
	} else {
		singularPattern = pattern[0:1]
	}
	isMatch := matchSingularPattern(line, singularPattern)
	if isMatch {
		return matchAtCurrentPoint(line[1:], pattern[len(singularPattern):])
	} else {
		return false
	}

}

func matchSingularPattern(line []byte, pattern string) bool {
	var characterSet string

	isSpecialPattern := contains(specialPatterns, pattern)
	isCharacterGroup := strings.HasPrefix(pattern, "[") && strings.HasSuffix(pattern, "]")
	isNegativeCharacterGroup := strings.HasPrefix(pattern, "[^") && strings.HasSuffix(pattern, "]")

	if isNegativeCharacterGroup {
		characterSet = pattern[2 : len(pattern)-1]
	} else if isCharacterGroup {
		characterSet = pattern[1 : len(pattern)-1]
	} else if isSpecialPattern {
		characterSet = getSpecialPatternCharacterSet(pattern)
	} else {
		characterSet = pattern[0:1]
	}

	isMatch := strings.ContainsRune(characterSet, rune(line[0]))
	fmt.Println(characterSet, string(line[0]), isMatch)
	if isNegativeCharacterGroup {
		return !isMatch
	} else {
		return isMatch
	}
}
