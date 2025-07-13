package ast

import "grpgscript/lex"

type Expr interface {}

type Binary struct {
    Left Expr
    Operator lex.Token
    Right Expr
}

type Grouping struct {
    Expression Expr
}

type Literal struct {
    Value any
}

type Unary struct {
    Operator lex.Token
    Right Expr
}
