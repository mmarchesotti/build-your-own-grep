package lexer

import (
	"strings"

	"github.com/mmarchesotti/build-your-own-grep/app/token"
)

func Parse(inputToken string) []token.Token {
	var tokens []token.Token

	for inputIndex := 0; inputIndex < len(inputToken); inputIndex++ {
		currentCharacter := inputToken[inputIndex]

		switch currentCharacter {
		case '\\':
			if inputIndex+1 < len(inputToken) {
				nextCharacter := inputToken[inputIndex+1]
				switch nextCharacter {
				case 'd':
					tokens = append(tokens, &token.Digit{})
					inputIndex += 1
				case 'w':
					tokens = append(tokens, &token.AlphaNumeric{})
					inputIndex += 1
				default:
					tokens = append(tokens, &token.Literal{Literal: '\\'})
				}
			} else {
				tokens = append(tokens, &token.Literal{Literal: '\\'})
			}
		case '[':
			closingIndex := strings.Index(inputToken[inputIndex:], "]")
			if closingIndex == -1 {
				tokens = append(tokens, &token.Literal{Literal: '['})
				continue
			}

			var groupTokens []token.Token
			groupCharacters := inputToken[inputIndex+1 : closingIndex]
			for groupIndex := 0; groupIndex < len(groupCharacters); groupIndex++ {
				currentGroupCharacter := groupCharacters[groupIndex]

				if currentGroupCharacter == '\\' && groupIndex+1 < len(inputToken) {
					nextCharacter := inputToken[inputIndex+1]
					switch nextCharacter {
					case 'd':
						groupTokens = append(groupTokens, &token.Digit{})
						groupIndex += 1
					case 'w':
						groupTokens = append(groupTokens, &token.AlphaNumeric{})
						groupIndex += 1
					default:
						groupTokens = append(groupTokens, &token.Literal{Literal: '\\'})
					}
				} else {
					groupTokens = append(groupTokens, &token.Literal{
						Literal: rune(groupCharacters[groupIndex]),
					})
				}
			}

			groupFirstCharacter := groupCharacters[0]
			if groupFirstCharacter == '^' {
				tokens = append(tokens, &token.NegativeGroup{Tokens: groupTokens[1:]})
			} else {
				tokens = append(tokens, &token.PositiveGroup{Tokens: groupTokens})
			}

			inputIndex += closingIndex + 1
		case '^':
			tokens = append(tokens, &token.StartAnchor{})
		case '$':
			tokens = append(tokens, &token.EndAnchor{})
		case '+':
			if len(tokens) > 0 {
				lastToken := tokens[len(tokens)-1]
				tokens[len(tokens)-1] = &token.PositiveClosure{SubToken: lastToken}
			} else {
				tokens = append(tokens, &token.Literal{Literal: '+'})
			}
		default:
			tokens = append(tokens, &token.Literal{
				Literal: rune(inputToken[inputIndex]),
			})
		}
	}

	return tokens
}
