package pattern

import (
	"unicode/utf8"

	"github.com/mmarchesotti/build-your-own-grep/app/engine"
	"github.com/mmarchesotti/build-your-own-grep/app/utils"
)

type PatternType string

const (
	INVALID        PatternType = "INVALID"
	SERIAL         PatternType = "SERIAL"
	LITERAL        PatternType = "LITERAL"
	DIGIT          PatternType = "DIGIT"
	ALPHANUMERIC   PatternType = "ALPHANUMERIC"
	POSITIVE_GROUP PatternType = "POSITIVE_GROUP"
	NEGATIVE_GROUP PatternType = "NEGATIVE_GROUP"
	START_ANCHOR   PatternType = "START_ANCHOR"
	END_ANCHOR     PatternType = "END_ANCHOR"
)

type basePattern struct {
	pType PatternType
}

func (pattern *basePattern) Type() string {
	return string(pattern.pType)
}

type Pattern interface {
	Type() string
	// Returns (bytesConsumed, didMatch)
	Match(ctx *engine.MatchContext) (int, bool)
}

type Literal struct {
	basePattern
	Literal rune
}

func (pattern *Literal) Match(ctx *engine.MatchContext) (int, bool) {
	return utf8.RuneLen(pattern.Literal),
		utils.StartsWithRune(ctx.RemainingLine(), pattern.Literal)
}

/* Digit Pattern */
type Digit struct {
	basePattern
}

func (pattern *Digit) Match(ctx *engine.MatchContext) (int, bool) {
	r, runeLen := ctx.FirstRune()
	if utils.IsDigit(r) {
		return runeLen, true
	} else {
		return 0, false
	}
}

/* Alphanumeric */
type AlphaNumeric struct {
	basePattern
}

func (pattern *AlphaNumeric) Match(ctx *engine.MatchContext) (int, bool) {
	r, runeLen := ctx.FirstRune()
	if utils.IsAlphaNumeric(r) {
		return runeLen, true
	} else {
		return 0, false
	}
}

/* Positive Groupattern */
type PositiveGroup struct {
	basePattern
	Patterns []Pattern
}

func (pattern *PositiveGroup) Match(ctx *engine.MatchContext) (int, bool) {
	for _, pattern := range pattern.Patterns {
		bytesConsumed, didMatch := pattern.Match(ctx)
		if didMatch {
			return bytesConsumed, true
		}
	}
	return 0, false
}

/* Negative groupattern */
type NegativeGroup struct {
	basePattern
	Patterns []Pattern
}

func (pattern *NegativeGroup) Match(ctx *engine.MatchContext) (int, bool) {
	for _, pattern := range pattern.Patterns {
		_, didMatch := pattern.Match(ctx)
		if didMatch {
			return 0, false
		}
	}
	_, firstRuneLen := ctx.FirstRune()
	return firstRuneLen, true
}

/* StartAnchor */
type StartAnchor struct {
	basePattern
}

func (pattern *StartAnchor) Match(ctx *engine.MatchContext) (int, bool) {
	return 0, ctx.IsAtStart
}

/* EndAnchor */
type EndAnchor struct {
	basePattern
}

func (pattern *EndAnchor) Match(ctx *engine.MatchContext) (int, bool) {
	return 0, ctx.IsAtEnd
}
