package parser

import (
	"fmt"

	"github.com/mmarchesotti/build-your-own-grep/app/ast"
	"github.com/mmarchesotti/build-your-own-grep/app/lexer"
	"github.com/mmarchesotti/build-your-own-grep/app/token"
)

type Parser struct {
	tokens   []token.Token
	position int
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{
		tokens:   tokens,
		position: 0,
	}
}

func (p *Parser) currentToken() token.Token {
	if p.position >= len(p.tokens) {
		return nil
	}
	return p.tokens[p.position]
}

func (p *Parser) consumeToken() token.Token {
	token := p.currentToken()
	p.position++
	return token
}

func (p *Parser) parseExpression() (ast.ASTNode, error) {
	node, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	for token.IsAlternation(p.currentToken()) {
		p.consumeToken()
		rightNode, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		node = &ast.AlternationNode{Left: node, Right: rightNode}
	}

	return node, nil
}

func (p *Parser) parseTerm() (ast.ASTNode, error) {
	node, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for token.CanConcatenate(p.currentToken()) {
		rightNode, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		node = &ast.ConcatenationNode{Left: node, Right: rightNode}
	}

	return node, nil
}

func (p *Parser) parseFactor() (ast.ASTNode, error) {
	node, err := p.parseAtom()
	if err != nil {
		return nil, err
	}

	for token.IsUnaryOperator(p.currentToken()) {
		t := p.consumeToken()
		switch t.(type) {
		case *token.OptionalQuantifier:
			node = &ast.OptionalNode{
				Child: node,
			}
		case *token.KleeneClosure:
			node = &ast.KleeneClosureNode{
				Child: node,
			}
		case *token.PositiveClosure:
			node = &ast.PositiveClosureNode{
				Child: node,
			}
		}
	}

	return node, nil
}

func (p *Parser) parseAtom() (ast.ASTNode, error) {
	switch t := p.currentToken().(type) {
	case *token.GroupingOpener:
		p.consumeToken()
		node, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		if !token.IsGroupingCloser(p.currentToken()) {
			return nil, fmt.Errorf("unmatched group opener")
		}
		p.consumeToken()
		return node, nil
	case *token.Literal:
		p.consumeToken()
		node := &ast.LiteralNode{
			Literal: t.Literal,
		}
		return node, nil
	case *token.CharacterSet:
		p.consumeToken()
		node := &ast.CharacterSetNode{
			IsPositive:       t.IsPositive,
			Literals:         t.Literals,
			Ranges:           t.Ranges,
			CharacterClasses: t.CharacterClasses,
		}
		return node, nil
	case *token.Wildcard:
		p.consumeToken()
		node := &ast.WildcardNode{}
		return node, nil
	case *token.Digit:
		p.consumeToken()
		node := &ast.DigitNode{}
		return node, nil
	case *token.AlphaNumeric:
		p.consumeToken()
		node := &ast.AlphaNumericNode{}
		return node, nil
	case *token.StartAnchor:
		p.consumeToken()
		node := &ast.StartAnchorNode{}
		return node, nil
	case *token.EndAnchor:
		p.consumeToken()
		node := &ast.EndAnchorNode{}
		return node, nil
	case *token.GroupingCloser:
		return nil, fmt.Errorf("unmatched group closer")
	default:
		return nil, fmt.Errorf("unexpected token: %T", t)
	}
}

func Parse(inputPattern string) (ast.ASTNode, error) {
	tokens, err := lexer.Tokenize(inputPattern)
	if err != nil {
		return nil, err
	}
	parser := NewParser(tokens)
	return parser.parseExpression()
}
