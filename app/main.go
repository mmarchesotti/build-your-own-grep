package main

import (
	"fmt"
	"io"
	"os"
	"unicode/utf8"

	"github.com/mmarchesotti/build-your-own-grep/app/buildnfa"
	"github.com/mmarchesotti/build-your-own-grep/app/nfa"
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

	ok := matchLine(line, pattern)

	if !ok {
		fmt.Println("Program returned 1")
		os.Exit(1)
	}

	fmt.Println("Program returned 0")
	// default exit code is 0 which means success
}

func matchLine(line []byte, inputPattern string) bool {
	nfa := buildnfa.Build(inputPattern)
	startNode := nfa.Start

	for lineStartIndex := 0; lineStartIndex <= len(line); {
		matchSucceeded := matchLineAt(startNode, line, lineStartIndex)
		if matchSucceeded {
			return true
		}

		_, size := utf8.DecodeRune(line[lineStartIndex:])
		if size == 0 {
			break
		}
		lineStartIndex += size
	}

	return false
}

func matchLineAt(n nfa.State, line []byte, lineIndex int) bool {
	switch node := n.(type) {
	case *nfa.MatcherState:
		if lineIndex >= len(line) {
			return false
		}

		nextRune, size := utf8.DecodeRune(line[lineIndex:])
		if node.Matcher.Match(nextRune) {
			lineIndex += size
			return matchLineAt(node.Out, line, lineIndex)
		} else {
			return false
		}
	case *nfa.SplitState:
		return matchLineAt(node.Branch1, line, lineIndex) ||
			matchLineAt(node.Branch2, line, lineIndex)
	case *nfa.AcceptingState:
		return true
	case *nfa.StartAnchorState:
		if lineIndex == 0 {
			return matchLineAt(node.Out, line, lineIndex)
		} else {
			return false
		}
	case *nfa.EndAnchorState:
		if lineIndex == len(line) {
			return matchLineAt(node.Out, line, lineIndex)
		} else {
			return false
		}
	}
	return true
}
