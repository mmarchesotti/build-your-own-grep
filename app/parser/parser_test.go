package parser

import (
	"reflect"
	"testing"

	"github.com/mmarchesotti/build-your-own-grep/app/ast"
)

// --- Test Helper Functions ---
// These make the test cases much cleaner and easier to read.
func lit(char rune) ast.ASTNode { return &ast.LiteralNode{Literal: char} }
func alt(left, right ast.ASTNode) ast.ASTNode {
	return &ast.AlternationNode{Left: left, Right: right}
}
func concat(left, right ast.ASTNode) ast.ASTNode {
	return &ast.ConcatenationNode{Left: left, Right: right}
}
func star(child ast.ASTNode) ast.ASTNode { return &ast.KleeneClosureNode{Child: child} }
func plus(child ast.ASTNode) ast.ASTNode { return &ast.PositiveClosureNode{Child: child} }
func opt(child ast.ASTNode) ast.ASTNode  { return &ast.OptionalNode{Child: child} }
func cs(neg bool, lits []rune) ast.ASTNode {
	return &ast.CharacterSetNode{Negated: neg, Literals: lits}
}

// --- Main Test Function ---

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ast.ASTNode
	}{
		{
			name:     "single literal",
			input:    "a",
			expected: lit('a'),
		},
		{
			name:     "simple concatenation",
			input:    "ab",
			expected: concat(lit('a'), lit('b')),
		},
		{
			name:     "long concatenation (will fail until bug is fixed)",
			input:    "abc",
			expected: concat(concat(lit('a'), lit('b')), lit('c')),
		},
		{
			name:     "simple alternation",
			input:    "a|b",
			expected: alt(lit('a'), lit('b')),
		},
		{
			name:     "alternation and concatenation precedence",
			input:    "ab|c",
			expected: alt(concat(lit('a'), lit('b')), lit('c')),
		},
		{
			name:     "parentheses for scope",
			input:    "a(b|c)",
			expected: concat(lit('a'), alt(lit('b'), lit('c'))),
		},
		{
			name:     "simple kleene star",
			input:    "a*",
			expected: star(lit('a')),
		},
		{
			name:     "kleene star on a group",
			input:    "(ab)*",
			expected: star(concat(lit('a'), lit('b'))),
		},
		{
			name:     "all quantifiers",
			input:    "a*b+c?",
			expected: concat(concat(star(lit('a')), plus(lit('b'))), opt(lit('c'))),
		},
		{
			name:     "character set",
			input:    "[abc]",
			expected: cs(false, []rune{'a', 'b', 'c'}),
		},
		{
			name:  "complex expression",
			input: "a(b|c)*d",
			expected: concat(
				concat(
					lit('a'),
					star(alt(lit('b'), lit('c'))),
				),
				lit('d'),
			),
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
