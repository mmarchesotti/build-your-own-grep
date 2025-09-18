package main

import (
	"strings"

	"github.com/codecrafters-io/grep-starter-go/app/pattern"
)

func Parse(inputPattern string) []pattern.Pattern {
	var patterns []pattern.Pattern

	for i := 0; i < len(inputPattern); i++ {
		currentCharacter := inputPattern[i]

		switch currentCharacter {
		case '\\':
			nextCharacter := inputPattern[i+1]
			switch nextCharacter {
			case 'd':
				patterns = append(patterns, &pattern.Digit{})
				i += 1
			case 'w':
				patterns = append(patterns, &pattern.AlphaNumeric{})
				i += 1
			default:
				patterns = append(patterns, &pattern.Literal{Literal: "\\"})
			}
		case '[':
			closingIndex := strings.Index(inputPattern[i+1:], "]")
			if closingIndex == -1 {
				patterns = append(patterns, &pattern.Literal{Literal: "["})
				continue
			}

			matches := 
			nextCharacter := inputPattern[i+1]
			
		default:
			patterns = append(patterns, &pattern.Literal{Literal: inputPattern[i : i+1]})
		}
	}

	return patterns
}


func Parse(inputPattern string) []pattern.Pattern {
	var patterns []pattern.Pattern

	for i := 0; i < len(inputPattern); i++ {
		currentCharacter := inputPattern[i]

		switch currentCharacter {
		case '\\':
			nextCharacter := inputPattern[i+1]
			switch nextCharacter {
			case 'd':
				patterns = append(patterns, &pattern.Digit{})
				i += 1
			case 'w':
				patterns = append(patterns, &pattern.AlphaNumeric{})
				i += 1
			default:
				patterns = append(patterns, &pattern.Literal{Literal: "\\"})
			}
		case '[':
			closingIndex := strings.Index(inputPattern[i+1:], "]")
			if closingIndex == -1 {
				patterns = append(patterns, &pattern.Literal{Literal: "["})
				continue
			}

			matches := 
			nextCharacter := inputPattern[i+1]
			
		default:
			patterns = append(patterns, &pattern.Literal{Literal: inputPattern[i : i+1]})
		}
	}

	return patterns
}
