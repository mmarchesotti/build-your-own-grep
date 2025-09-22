package token

import (
	"unicode/utf8"

	"github.com/mmarchesotti/build-your-own-grep/app/engine"
	"github.com/mmarchesotti/build-your-own-grep/app/utils"
)

type TokenType string

const (
	INVALID             TokenType = "INVALID"
	SERIAL              TokenType = "SERIAL"
	LITERAL             TokenType = "LITERAL"
	DIGIT               TokenType = "DIGIT"
	ALPHANUMERIC        TokenType = "ALPHANUMERIC"
	POSITIVE_GROUP      TokenType = "POSITIVE_GROUP"
	NEGATIVE_GROUP      TokenType = "NEGATIVE_GROUP"
	START_ANCHOR        TokenType = "START_ANCHOR"
	END_ANCHOR          TokenType = "END_ANCHOR"
	KLEENE_CLOSURE      TokenType = "KLEENE_CLOSURE"
	POSITIVE_CLOSURE    TokenType = "POSITIVE_CLOSURE"
	OPTIONAL_QUANTIFIER TokenType = "OPTIONAL_QUANTIFIER"
	WILDCARD            TokenType = "WILDCARD"
)

type baseToken struct {
	pType TokenType
}

func (token *baseToken) Type() string {
	return string(token.pType)
}

type Token interface {
	Type() string
	// Returns (bytesConsumed, didMatch)
	Match(ctx *engine.MatchContext) (int, bool)
}

type Literal struct {
	baseToken
	Literal rune
}

func (token *Literal) Match(ctx *engine.MatchContext) (int, bool) {
	return utf8.RuneLen(token.Literal),
		utils.StartsWithRune(ctx.RemainingLine(), token.Literal)
}

/* Digit Token */
type Digit struct {
	baseToken
}

func (token *Digit) Match(ctx *engine.MatchContext) (int, bool) {
	r, runeLen := ctx.FirstRune()
	if utils.IsDigit(r) {
		return runeLen, true
	} else {
		return 0, false
	}
}

/* Alphanumeric */
type AlphaNumeric struct {
	baseToken
}

func (token *AlphaNumeric) Match(ctx *engine.MatchContext) (int, bool) {
	r, runeLen := ctx.FirstRune()
	if utils.IsAlphaNumeric(r) {
		return runeLen, true
	} else {
		return 0, false
	}
}

/* Positive Group token */
type PositiveGroup struct {
	baseToken
	Tokens []Token
}

func (token *PositiveGroup) Match(ctx *engine.MatchContext) (int, bool) {
	for _, token := range token.Tokens {
		bytesConsumed, didMatch := token.Match(ctx)
		if didMatch {
			return bytesConsumed, true
		}
	}
	return 0, false
}

/* Negative group token */
type NegativeGroup struct {
	baseToken
	Tokens []Token
}

func (token *NegativeGroup) Match(ctx *engine.MatchContext) (int, bool) {
	for _, token := range token.Tokens {
		_, didMatch := token.Match(ctx)
		if didMatch {
			return 0, false
		}
	}
	_, firstRuneLen := ctx.FirstRune()
	return firstRuneLen, true
}

/* StartAnchor */
type StartAnchor struct {
	baseToken
}

func (token *StartAnchor) Match(ctx *engine.MatchContext) (int, bool) {
	return 0, ctx.IsAtStart
}

/* EndAnchor */
type EndAnchor struct {
	baseToken
}

func (token *EndAnchor) Match(ctx *engine.MatchContext) (int, bool) {
	return 0, ctx.IsAtEnd
}

/* PositiveClosure */
type PositiveClosure struct {
	baseToken
	SubToken Token
}

func (token *PositiveClosure) Match(ctx *engine.MatchContext) (int, bool) {
	bytesConsumed, didMatch := token.SubToken.Match(ctx)
	if !didMatch {
		return 0, false
	}

	totalBytesConsumed := bytesConsumed
	localCtx := *ctx
	localCtx.Index += bytesConsumed
	for {
		bytesConsumed, didMatch := token.SubToken.Match(&localCtx)
		if !didMatch {
			break
		}

		totalBytesConsumed += bytesConsumed
		localCtx.Index += bytesConsumed
	}

	return totalBytesConsumed, true
}
