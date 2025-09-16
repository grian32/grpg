package ast

import (
	"bytes"
	"grpgscript/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
	Pos() Position
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

func (p *Program) Pos() Position {
	first := p.Statements[0].Pos()
	last := p.Statements[len(p.Statements)-1].Pos()

	return Position{
		Start:     first.Start,
		End:       last.End,
		StartLine: first.StartLine,
		EndLine:   last.EndLine,
	}
}
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type VarStatement struct {
	Token token.Token // token.let
	Name  *Identifier
	Value Expression
}

func (vs *VarStatement) Pos() Position {
	return Position{
		Start:     vs.Token.Col,
		End:       vs.Value.Pos().End,
		StartLine: vs.Token.Line,
		EndLine:   vs.Value.Pos().EndLine,
	}
}
func (vs *VarStatement) statementNode()       { /* noop */ }
func (vs *VarStatement) TokenLiteral() string { return vs.Token.Literal }
func (vs *VarStatement) String() string {
	var out bytes.Buffer

	out.WriteString(vs.TokenLiteral() + " " + vs.Name.String() + " = ")

	if vs.Value != nil {
		out.WriteString(vs.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type Identifier struct {
	Token token.Token // token.ident
	Value string
}

func (i *Identifier) Pos() Position {
	return Position{
		Start:     i.Token.Col,
		End:       i.Token.End,
		StartLine: i.Token.Line,
		EndLine:   i.Token.Line,
	}
}
func (i *Identifier) expressionNode()      { /* noop */ }
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) Pos() Position {
	return Position{
		Start:     rs.Token.Col,
		End:       rs.ReturnValue.Pos().End,
		StartLine: rs.Token.Line,
		EndLine:   rs.ReturnValue.Pos().EndLine,
	}
}
func (rs *ReturnStatement) statementNode()       { /* noop */ }
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) Pos() Position {
	return es.Expression.Pos()
}
func (es *ExpressionStatement) statementNode()       { /* noop */ }
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) Pos() Position {
	return Position{
		Start:     il.Token.Col,
		End:       il.Token.End,
		StartLine: il.Token.Line,
		EndLine:   il.Token.Line,
	}
}
func (il *IntegerLiteral) expressionNode()      { /* noop */ }
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) Pos() Position {
	return Position{
		Start:     pe.Token.Col,
		End:       pe.Right.Pos().End,
		StartLine: pe.Token.Line,
		EndLine:   pe.Token.Line,
	}
}
func (pe *PrefixExpression) expressionNode()      { /* noop */ }
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	return "(" + pe.Operator + pe.Right.String() + ")"
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) Pos() Position {
	return Position{
		Start:     ie.Left.Pos().Start,
		End:       ie.Right.Pos().End,
		StartLine: ie.Left.Pos().StartLine,
		EndLine:   ie.Right.Pos().EndLine,
	}
}
func (ie *InfixExpression) expressionNode()      { /* noop */ }
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	return "(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")"
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) Pos() Position {
	return Position{
		Start:     b.Token.Col,
		End:       b.Token.End,
		StartLine: b.Token.Line,
		EndLine:   b.Token.Line,
	}
}
func (b *Boolean) expressionNode()      { /* noop */ }
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) Pos() Position {
	var endPos uint64 = 0
	var endLine uint64 = 0

	if ie.Alternative != nil {
		endPos = ie.Alternative.Pos().End
		endLine = ie.Alternative.Pos().End
	} else {
		endPos = ie.Consequence.Pos().End
		endLine = ie.Consequence.Pos().End
	}

	return Position{
		Start:     ie.Token.Col,
		End:       endPos,
		StartLine: ie.Token.Line,
		EndLine:   endLine,
	}
}
func (ie *IfExpression) expressionNode()      { /* noop */ }
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if" + ie.Condition.String() + " " + ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else " + ie.Alternative.String())
	}

	return out.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) Pos() Position {
	first := bs.Statements[0]
	last := bs.Statements[len(bs.Statements)-1]
	return Position{
		Start:     first.Pos().Start,
		End:       last.Pos().End,
		StartLine: first.Pos().StartLine,
		EndLine:   first.Pos().EndLine,
	}
}
func (bs *BlockStatement) statementNode()       { /* noop */ }
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) Pos() Position {
	return Position{
		Start:     fl.Token.Col,
		End:       fl.Body.Pos().End,
		StartLine: fl.Token.Line,
		EndLine:   fl.Body.Pos().EndLine,
	}
}
func (fl *FunctionLiteral) expressionNode()      { /* noop */ }
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral() + "(" + strings.Join(params, ", ") + ") " + fl.Body.String())

	return out.String()
}

type CallExpresion struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpresion) Pos() Position {
	if len(ce.Arguments) > 0 {
		last := ce.Arguments[len(ce.Arguments)-1]

		return Position{
			Start:     ce.Token.Col,
			End:       last.Pos().End,
			StartLine: ce.Token.Line,
			EndLine:   last.Pos().EndLine,
		}
	}

	return Position{
		Start:     ce.Token.Col,
		End:       ce.Token.End,
		StartLine: ce.Token.Line,
		EndLine:   ce.Token.Line,
	}
}
func (ce *CallExpresion) expressionNode()      { /* noop */ }
func (ce *CallExpresion) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpresion) String() string {
	var out bytes.Buffer

	args := []string{}

	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String() + "(" + strings.Join(args, ", ") + ")")

	return out.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) Pos() Position {
	return Position{
		Start:     sl.Token.Col,
		End:       sl.Token.End,
		StartLine: sl.Token.Line,
		EndLine:   sl.Token.Line,
	}
}
func (sl *StringLiteral) expressionNode()      { /* noop */ }
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) Pos() Position {
	last := al.Elements[len(al.Elements)-1]

	return Position{
		Start:     al.Token.Col,
		End:       last.Pos().End,
		StartLine: al.Token.Line,
		EndLine:   last.Pos().EndLine,
	}
}
func (al *ArrayLiteral) expressionNode()      { /* noop */ }
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	return "[" + strings.Join(elements, ", ") + "]"
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) Pos() Position {
	return Position{
		Start:     ie.Token.Col,
		End:       ie.Index.Pos().End,
		StartLine: ie.Token.Line,
		EndLine:   ie.Index.Pos().EndLine,
	}
}
func (ie *IndexExpression) expressionNode()      { /* noop */ }
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	return "(" + ie.Left.String() + "[" + ie.Index.String() + "])"
}

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) Pos() Position {
	var last Expression

	tokenPos := Position{
		Start:     hl.Token.Col,
		End:       hl.Token.End,
		StartLine: hl.Token.Line,
		EndLine:   hl.Token.Line,
	}

	if hl.Pairs == nil {
		return tokenPos
	}

	for _, expr := range hl.Pairs {
		if last == nil {
			last = expr
			continue
		}

		if expr.Pos().EndLine > last.Pos().EndLine || expr.Pos().End > last.Pos().End {
			last = expr
		}
	}

	if last == nil {
		return tokenPos
	}

	return Position{
		Start:     hl.Token.Col,
		End:       last.Pos().End,
		StartLine: hl.Token.Line,
		EndLine:   last.Pos().EndLine,
	}
}
func (hl *HashLiteral) expressionNode()      { /* noop */ }
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	return "{" + strings.Join(pairs, "") + "}"
}
