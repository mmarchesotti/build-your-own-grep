package parser

import (
	"github.com/mmarchesotti/build-your-own-grep/app/lexer"
	"github.com/mmarchesotti/build-your-own-grep/app/token"
)

func insertConcatenations(tokens []token.Token) []token.Token {
	if len(tokens) < 2 {
		return tokens
	}

	tokensWithConcatenations := make([]token.Token, 0, len(tokens)*2)
	tokensWithConcatenations = append(tokensWithConcatenations, tokens[0])

	for i := 1; i < len(tokens); i++ {
		prev := tokens[i-1]
		curr := tokens[i]

		if token.IsEnder(prev) && token.IsStarter(curr) {
			tokensWithConcatenations = append(tokensWithConcatenations, &token.Concatenation{})
		}
		tokensWithConcatenations = append(tokensWithConcatenations, curr)
	}

	return tokensWithConcatenations
}

func Parse(inputPattern string) []token.Token {
	tokens := lexer.Tokenize(inputPattern)
	tokensWithConcatenations := insertConcatenations(tokens)

	postOrderedTokens := make([]token.Token, 0, len(tokensWithConcatenations))
	waitingOperators := make([]token.Token, 0, len(tokensWithConcatenations))

	for _, currentToken := range tokensWithConcatenations {
		if token.IsOperator(currentToken) {
			for len(waitingOperators) > 0 {
				topOperator := waitingOperators[len(waitingOperators)-1]
				if token.Precedence(currentToken) > token.Precedence(topOperator) {
					break
				}
				postOrderedTokens = append(postOrderedTokens, topOperator)
				waitingOperators = waitingOperators[:len(waitingOperators)-1]
			}
			waitingOperators = append(waitingOperators, currentToken)
		} else {
			postOrderedTokens = append(postOrderedTokens, currentToken)
		}
	}

	for len(waitingOperators) > 0 {
		postOrderedTokens = append(postOrderedTokens, waitingOperators[len(waitingOperators)-1])
		waitingOperators = waitingOperators[:len(waitingOperators)-1]
	}

	return postOrderedTokens
}
