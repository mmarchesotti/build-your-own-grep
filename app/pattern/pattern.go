package pattern

import (
	"bytes"
	"slices"

	"github.com/codecrafters-io/grep-starter-go/app/utils"
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

type Pattern interface {
	Type() PatternType
	Match([]byte) ([]byte, bool)
}

/* Literal Pattern */
type Literal struct {
	Literal string
}

func (p *Literal) Type() PatternType {
	return LITERAL
}

func (p *Literal) Match(r []byte) ([]byte, bool) {
	i := bytes.Index(r, []byte(p.Literal))
	if i == -1 {
		return r, false
	}
	return r[i:], true
}

/* Digit Pattern */
type Digit struct{}

func (p *Digit) Type() PatternType {
	return DIGIT
}

func (p *Digit) Match(r []byte) ([]byte, bool) {
	for i := range r {
		if utils.IsDigit(rune(r[i])) {
			return r[i+1:], true
		}
	}
	return r, false
}

/* Alphanumeric */
type AlphaNumeric struct{}

func (p *AlphaNumeric) Type() PatternType {
	return ALPHANUMERIC
}

func (p *AlphaNumeric) Match(r []byte) ([]byte, bool) {
	for i := range r {
		if utils.IsAlphaNumeric(rune(r[i])) {
			return r[i+1:], true
		}
	}
	return r, false
}

/* Positive Group */
type PositiveGroup struct {
	Patterns []Pattern
}

func (p *PositiveGroup) Type() PatternType {
	return POSITIVE_GROUP
}

func (p *PositiveGroup) Match(r []byte) ([]byte, bool) {
	for i := range r {
		if slices.Contains(p.Matches, r[i]) {
			return r[i+1:], true
		}
	}
	return r, false
}

/* Negative group */
type NegativeGroup struct {
	Matches []byte
}

func (p *NegativeGroup) Type() PatternType {
	return NEGATIVE_GROUP
}

func (p *NegativeGroup) Match(r []byte) ([]byte, bool) {
	for i := range r {
		if !slices.Contains(p.Matches, r[i]) {
			return r[i+1:], true
		}
	}
	return r, false
}

/* StartAnchor */
type StartAnchor struct {
	Matches []byte
}

func (p *StartAnchor) Type() PatternType {
	return START_ANCHOR
}

func (p *StartAnchor) Match(r []byte) ([]byte, bool) {
	patternLen := len(p.Matches)
	rLen := len(r)

	if rLen < patternLen {
		return r, false
	}

	if bytes.Equal(p.Matches, r[:rLen]) {
		return r[rLen:], true
	}

	return r, false
}

/* EndAnchor */
type EndAnchor struct{}

func (p *EndAnchor) Type() PatternType {
	return END_ANCHOR
}

func (p *EndAnchor) Match(r []byte) ([]byte, bool) {
	patternLen := len(p.Matches)
	rLen := len(r)

	if rLen < patternLen {
		return r, false
	}

	if bytes.Equal(p.Matches, r[rLen-patternLen:]) {
		return r[rLen:], true
	}

	return r, false
}
