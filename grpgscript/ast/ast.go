package ast

import "grpgscript/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type VarStatement struct {
	Token token.Token // token.let
	Name *Identifier
	Value Expression
}

func (ls *VarStatement) statementNode() { /* noop */}
func (ls *VarStatement) TokenLiteral() string { return ls.Token.Literal }

type Identifier struct {
	Token token.Token // token.ident
	Value string
}

func (i *Identifier) statementNode() { /* noop */}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }


type ReturnStatement struct {
	Token token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() { /* noop */ }
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
