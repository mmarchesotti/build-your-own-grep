package engine

import "unicode/utf8"

type MatchContext struct {
	Line      []byte
	Index     int
	IsAtStart bool
	IsAtEnd   bool
}

func (ctx *MatchContext) RemainingLine() []byte {
	return ctx.Line[ctx.Index:]
}

func (ctx *MatchContext) FirstRune() (rune, int) {
	return utf8.DecodeRune(ctx.RemainingLine())
}
