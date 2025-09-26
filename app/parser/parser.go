package parser

import (
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

func (p *Parser) parseExpression() ast.ASTNode {
	node := p.parseTerm()

	for token.IsAlternation(p.currentToken()) {
		p.consumeToken()
		rightNode := p.parseTerm()
		node = &ast.AlternationNode{Left: node, Right: rightNode}
	}

	return node
}

func (p *Parser) parseTerm() ast.ASTNode {
	node := p.parseFactor()

	for token.CanConcatenate(p.currentToken()) {
		rightNode := p.parseFactor()
		node = &ast.ConcatenationNode{Left: node, Right: rightNode}
	}

	return node
}

func (p *Parser) parseFactor() ast.ASTNode {
	node := p.parseAtom()

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

	return node
}

func (p *Parser) parseAtom() ast.ASTNode {
	switch t := p.currentToken().(type) {
	case *token.GroupingOpener:
		p.consumeToken()
		node := p.parseExpression()

		if !token.IsGroupingCloser(p.currentToken()) {
			// TODO ERROR UNMATCHED GROUP OPENER
		}
		p.consumeToken()
		return node
	case *token.Literal:
		p.consumeToken()
		return &ast.LiteralNode{
			Literal: t.Literal,
		}
	case *token.CharacterSet:
		p.consumeToken()
		return &ast.CharacterSetNode{
			IsPositive:       t.IsPositive,
			Literals:         t.Literals,
			Ranges:           t.Ranges,
			CharacterClasses: t.CharacterClasses,
		}
	case *token.Wildcard:
		p.consumeToken()
		return &ast.WildcardNode{}
	case *token.Digit:
		p.consumeToken()
		return &ast.DigitNode{}
	case *token.AlphaNumeric:
		p.consumeToken()
		return &ast.AlphaNumericNode{}
	case *token.StartAnchor:
		p.consumeToken()
		return &ast.StartAnchorNode{}
	case *token.EndAnchor:
		p.consumeToken()
		return &ast.EndAnchorNode{}
	case *token.GroupingCloser:
		// TODO ERROR UNMATCHED GROUP CLOSER
		return nil
	default:
		// TODO ERROR UNEXPECTED TOKEN
		return nil
	}
}

func Parse(inputPattern string) ast.ASTNode {
	tokens := lexer.Tokenize(inputPattern)
	parser := NewParser(tokens)
	return parser.parseExpression()
}
