package lexer

import (
	"reflect"
	"testing"

	"github.com/mmarchesotti/build-your-own-grep/app/token"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []token.Token
	}{
		{
			name:  "simple literals",
			input: "abc",
			expected: []token.Token{
				&token.Literal{Literal: 'a'},
				&token.Literal{Literal: 'b'},
				&token.Literal{Literal: 'c'},
			},
		},
		{
			name:  "all metacharacters",
			input: `*+?|^$.`,
			expected: []token.Token{
				&token.KleeneClosure{},
				&token.PositiveClosure{},
				&token.OptionalQuantifier{},
				&token.Alternation{},
				&token.StartAnchor{},
				&token.EndAnchor{},
				&token.Wildcard{},
			},
		},
		{
			name:  "escaped metacharacter",
			input: `\+`,
			expected: []token.Token{
				&token.Literal{Literal: '+'},
			},
		},
		{
			name:  "escaped predefined classes",
			input: `\d\w`,
			expected: []token.Token{
				&token.Digit{},
				&token.AlphaNumeric{},
			},
		},
		{
			name:  "simple character set",
			input: "[abc]",
			expected: []token.Token{
				&token.CharacterSet{
					Literals: []rune{'a', 'b', 'c'},
				},
			},
		},
		{
			name:  "negated character set",
			input: "[^abc]",
			expected: []token.Token{
				&token.CharacterSet{
					Negated:  true,
					Literals: []rune{'a', 'b', 'c'},
				},
			},
		},
		{
			name:  "character set with escaped class",
			input: `[a\d]`,
			expected: []token.Token{
				&token.CharacterSet{
					Negated:          false,
					Literals:         []rune{'a'},
					CharacterClasses: []token.PredefinedClass{token.ClassDigit},
				},
			},
		},
		{
			name:  "empty character set",
			input: `[]`,
			expected: []token.Token{
				&token.CharacterSet{
					Negated: false,
				},
			},
		},
		{
			name:  "literal and character set concatenation",
			input: `a[bc]`,
			expected: []token.Token{
				&token.Literal{Literal: 'a'},
				&token.CharacterSet{
					Negated:  false,
					Literals: []rune{'b', 'c'},
				},
			},
		},
		// {
		// 	name:     "unmatched opening bracket",
		// 	input:    `[abc`,
		// 	expected: []token.Token{}, // TODO TEST ERROR HANDLING
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := Tokenize(tt.input)

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("Tokenize() for input '%s' failed", tt.input)
				t.Errorf("got:  %#v", actual)
				t.Errorf("want: %#v", tt.expected)
			}
		})
	}
}
