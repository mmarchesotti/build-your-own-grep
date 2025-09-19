package lexer

import (
	"strings"

	"github.com/mmarchesotti/build-your-own-grep/app/pattern"
)

func Parse(inputPattern string) []pattern.Pattern {
	var patterns []pattern.Pattern

	for inputIndex := 0; inputIndex < len(inputPattern); inputIndex++ {
		currentCharacter := inputPattern[inputIndex]

		switch currentCharacter {
		case '\\':
			if inputIndex+1 < len(inputPattern) {
				nextCharacter := inputPattern[inputIndex+1]
				switch nextCharacter {
				case 'd':
					patterns = append(patterns, &pattern.Digit{})
					inputIndex += 1
				case 'w':
					patterns = append(patterns, &pattern.AlphaNumeric{})
					inputIndex += 1
				default:
					patterns = append(patterns, &pattern.Literal{Literal: '\\'})
				}
			} else {
				patterns = append(patterns, &pattern.Literal{Literal: '\\'})
			}
		case '[':
			closingIndex := strings.Index(inputPattern[inputIndex:], "]")
			if closingIndex == -1 {
				patterns = append(patterns, &pattern.Literal{Literal: '['})
				continue
			}

			var groupPatterns []pattern.Pattern
			groupCharacters := inputPattern[inputIndex+1 : closingIndex]
			for groupIndex := 0; groupIndex < len(groupCharacters); groupIndex++ {
				currentGroupCharacter := groupCharacters[groupIndex]

				if currentGroupCharacter == '\\' && groupIndex+1 < len(inputPattern) {
					nextCharacter := inputPattern[inputIndex+1]
					switch nextCharacter {
					case 'd':
						groupPatterns = append(groupPatterns, &pattern.Digit{})
						groupIndex += 1
					case 'w':
						groupPatterns = append(groupPatterns, &pattern.AlphaNumeric{})
						groupIndex += 1
					default:
						groupPatterns = append(groupPatterns, &pattern.Literal{Literal: '\\'})
					}
				} else {
					groupPatterns = append(groupPatterns, &pattern.Literal{
						Literal: rune(groupCharacters[groupIndex]),
					})
				}
			}

			groupFirstCharacter := groupCharacters[0]
			if groupFirstCharacter == '^' {
				patterns = append(patterns, &pattern.NegativeGroup{Patterns: groupPatterns[1:]})
			} else {
				patterns = append(patterns, &pattern.PositiveGroup{Patterns: groupPatterns})
			}

			inputIndex += closingIndex + 1
		case '^':
			patterns = append(patterns, &pattern.StartAnchor{})
		case '$':
			patterns = append(patterns, &pattern.EndAnchor{})
		default:
			patterns = append(patterns, &pattern.Literal{
				Literal: rune(inputPattern[inputIndex]),
			})
		}
	}

	return patterns
}
