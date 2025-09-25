package matcher

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isLower(r rune) bool {
	return r >= 'a' && r <= 'z'
}

func isUpper(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

func isAlpha(r rune) bool {
	return isUpper(r) || isLower(r)
}

func isAlphaNumeric(r rune) bool {
	return isAlpha(r) || isDigit(r) || r == '_'
}

func match(r rune, rng [2]rune) bool {
	if rng[1] > rng[0] {
		// TODO ERROR RANGE VALUES REVERSED
	}
	return r >= rng[0] && r <= rng[1]
}

type Matcher interface {
	Match(r rune) bool
}

type PredefinedClassMatcher interface {
	Matcher
	isPredefinedClass()
}

type LiteralMatcher struct {
	Literal rune
}

func (l *LiteralMatcher) Match(r rune) bool {
	return r == l.Literal
}

type CharacterSetMatcher struct {
	IsPositive               bool
	Literals                 []rune
	Ranges                   [][2]rune
	CharacterClassesMatchers []PredefinedClassMatcher
}

func (p *CharacterSetMatcher) Match(r rune) bool {
	for _, literal := range p.Literals {
		if r == literal {
			return p.IsPositive
		}
	}
	for _, rng := range p.Ranges {
		if match(r, rng) {
			return p.IsPositive
		}
	}
	for _, characterClass := range p.CharacterClassesMatchers {
		if characterClass.Match(r) {
			return p.IsPositive
		}
	}
	return !p.IsPositive
}

type WildcardMatcher struct{}

func (w *WildcardMatcher) Match(r rune) bool {
	return true
}

type DigitMatcher struct{}

func (d *DigitMatcher) Match(r rune) bool {
	return isDigit(r)
}

func (d *DigitMatcher) isPredefinedClass() {}

type AlphaNumericMatcher struct{}

func (a *AlphaNumericMatcher) Match(r rune) bool {
	return isAlphaNumeric(r)
}

func (a *AlphaNumericMatcher) isPredefinedClass() {}
