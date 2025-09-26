package lexer

import (
	"reflect"
	"testing"

	"github.com/mmarchesotti/build-your-own-grep/app/predefinedclass"
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
					IsPositive: true,
					Literals:   []rune{'a', 'b', 'c'},
				},
			},
		},
		{
			name:  "negated character set",
			input: "[^abc]",
			expected: []token.Token{
				&token.CharacterSet{
					IsPositive: false,
					Literals:   []rune{'a', 'b', 'c'},
				},
			},
		},
		{
			name:  "character set with escaped class",
			input: `[a\d]`,
			expected: []token.Token{
				&token.CharacterSet{
					IsPositive: true,
					Literals:   []rune{'a'},
					CharacterClasses: []predefinedclass.PredefinedClass{
						predefinedclass.ClassDigit,
					},
				},
			},
		},
		{
			name:  "empty character set",
			input: `[]`,
			expected: []token.Token{
				&token.CharacterSet{
					IsPositive: true,
				},
			},
		},
		{
			name:  "literal and character set concatenation",
			input: `a[bc]`,
			expected: []token.Token{
				&token.Literal{Literal: 'a'},
				&token.CharacterSet{
					IsPositive: true,
					Literals:   []rune{'b', 'c'},
				},
			},
		},
		{
			name:  "combination of multiple characters",
			input: `a(b|c)*d`,
			expected: []token.Token{
				&token.Literal{Literal: 'a'},
				&token.GroupingOpener{},
				&token.Literal{Literal: 'b'},
				&token.Alternation{},
				&token.Literal{Literal: 'c'},
				&token.GroupingCloser{},
				&token.KleeneClosure{},
				&token.Literal{Literal: 'd'},
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
