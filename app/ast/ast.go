package ast

import predefinedclass "github.com/mmarchesotti/build-your-own-grep/app/predefinedclass"

type ASTNode interface {
	isASTNode()
}

type baseASTNode struct{}

func (n *baseASTNode) isASTNode() {}

type AlternationNode struct {
	baseASTNode
	Left  ASTNode
	Right ASTNode
}

type ConcatenationNode struct {
	baseASTNode
	Left  ASTNode
	Right ASTNode
}

type KleeneClosureNode struct {
	baseASTNode
	Child ASTNode
}

type PositiveClosureNode struct {
	baseASTNode
	Child ASTNode
}

type OptionalNode struct {
	baseASTNode
	Child ASTNode
}

type LiteralNode struct {
	baseASTNode
	Literal rune
}

type CharacterSetNode struct {
	baseASTNode
	Negated          bool
	Literals         []rune
	Ranges           [][2]rune
	CharacterClasses []predefinedclass.PredefinedClass
}

type WildcardNode struct {
	baseASTNode
}

type DigitNode struct {
	baseASTNode
}

type AlphaNumericNode struct {
	baseASTNode
}

type StartAnchorNode struct {
	baseASTNode
}

type EndAnchorNode struct {
	baseASTNode
}
