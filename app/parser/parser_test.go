package parser

import (
	"reflect"
	"testing"

	"github.com/mmarchesotti/build-your-own-grep/app/token"
)

// Helper to make test definitions cleaner
func l(char rune) *token.Literal          { return &token.Literal{Literal: char} }
func alt() *token.Alternation             { return &token.Alternation{} }
func concat() *token.Concatenation        { return &token.Concatenation{} }
func kleene() *token.KleeneClosure        { return &token.KleeneClosure{} }
func pos() *token.PositiveClosure         { return &token.PositiveClosure{} }
func opt() *token.OptionalQuantifier      { return &token.OptionalQuantifier{} }
func cs(lits ...rune) *token.CharacterSet { return &token.CharacterSet{Literals: lits} }

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []token.Token
	}{
		{
			name:  "simple concatenation",
			input: "ab",
			expected: []token.Token{
				l('a'), l('b'), concat(),
			},
		},
		{
			name:  "simple alternation",
			input: "a|b",
			expected: []token.Token{
				l('a'), l('b'), alt(),
			},
		},
		{
			name:  "alternation and concatenation precedence",
			input: "ab|c",
			expected: []token.Token{
				l('a'), l('b'), concat(), l('c'), alt(),
			},
		},
		{
			name:  "concatenation with quantifiers",
			input: "a*b",
			expected: []token.Token{
				l('a'), kleene(), l('b'), concat(),
			},
		},
		{
			name:  "complex expression with multiple operators",
			input: "a*|b+",
			expected: []token.Token{
				l('a'), kleene(), l('b'), pos(), alt(),
			},
		},
		{
			name:  "character set concatenation",
			input: "a[bc]",
			expected: []token.Token{
				l('a'), cs('b', 'c'), concat(),
			},
		},
		{
			name:  "quantifier on a character set",
			input: "[ab]?",
			expected: []token.Token{
				cs('a', 'b'), opt(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := Parse(tt.input)

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("Parse() for input '%s' failed", tt.input)
				t.Errorf("got:  %#v", actual)
				t.Errorf("want: %#v", tt.expected)
			}
		})
	}
}
