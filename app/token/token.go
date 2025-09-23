package token

type TokenType string

const (
	INVALID             TokenType = "INVALID"
	LITERAL             TokenType = "LITERAL"
	DIGIT               TokenType = "DIGIT"
	ALPHANUMERIC        TokenType = "ALPHANUMERIC"
	CHARACTER_SET       TokenType = "CHARACTER_SET"
	START_ANCHOR        TokenType = "START_ANCHOR"
	END_ANCHOR          TokenType = "END_ANCHOR"
	KLEENE_CLOSURE      TokenType = "KLEENE_CLOSURE"
	POSITIVE_CLOSURE    TokenType = "POSITIVE_CLOSURE"
	OPTIONAL_QUANTIFIER TokenType = "OPTIONAL_QUANTIFIER"
	WILDCARD            TokenType = "WILDCARD"
	ALTERNATION         TokenType = "ALTERNATION"
	CONCATENATION       TokenType = "CONCATENATION"
	GROUPING_OPENER     TokenType = "GROUPING_OPENER"
	GROUPING_CLOSER     TokenType = "GROUPING_CLOSER"
)

type PredefinedClass int

const (
	ClassDigit PredefinedClass = iota
	ClassAlphanumeric
	ClassWhitespace
)

// --- Helper Functions ---

func IsOperator(t Token) bool {
	switch t.(type) {
	case *KleeneClosure, *PositiveClosure, *OptionalQuantifier, *Alternation, *Concatenation:
		return true
	default:
		return false
	}
}

func IsStarter(t Token) bool {
	switch t.(type) {
	case *Literal, *CharacterSet, *Wildcard, *Digit, *AlphaNumeric, *StartAnchor:
		return true
	default:
		return false
	}
}

func IsEnder(t Token) bool {
	switch t.(type) {
	case *Literal, *CharacterSet, *Wildcard, *Digit, *AlphaNumeric, *EndAnchor,
		*KleeneClosure, *PositiveClosure, *OptionalQuantifier:
		return true
	default:
		return false
	}
}

func Precedence(t Token) int {
	switch t.(type) {
	case *KleeneClosure, *PositiveClosure, *OptionalQuantifier:
		return 3
	case *Concatenation:
		return 2
	case *Alternation:
		return 1
	default:
		return 0
	}
}

// --- Token Interface and Structs ---

type Token interface {
	Type() string
}

type baseToken struct {
	pType TokenType
}

func (token *baseToken) Type() string {
	return string(token.pType)
}

// OPERATORS
type Concatenation struct{ baseToken }
type KleeneClosure struct{ baseToken }
type PositiveClosure struct{ baseToken }
type OptionalQuantifier struct{ baseToken }
type Alternation struct{ baseToken }

// OPERANDS
type Literal struct {
	baseToken
	Literal rune
}
type Digit struct{ baseToken }
type AlphaNumeric struct{ baseToken }
type CharacterSet struct {
	baseToken
	Negated          bool
	Literals         []rune
	Ranges           [][2]rune
	CharacterClasses []PredefinedClass
}
type StartAnchor struct{ baseToken }
type EndAnchor struct{ baseToken }
type Wildcard struct{ baseToken }
type GroupingOpener struct{ baseToken }
type GroupingCloser struct{ baseToken }
