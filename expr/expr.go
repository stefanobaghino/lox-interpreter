package expr

import "lox/token"

type Visitor interface {
	VisitBinaryExpr(*Binary) interface{}
	VisitGroupingExpr(*Grouping) interface{}
	VisitLiteralExpr(*Literal) interface{}
	VisitUnaryExpr(*Unary) interface{}
}

func (e *Binary) Accept(v Visitor) interface{} {
	return v.VisitBinaryExpr(e)
}

type Expr interface {
	Accept(Visitor) interface{}
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

type Grouping struct {
	Expression Expr
}

func (e *Grouping) Accept(v Visitor) interface{} {
	return v.VisitGroupingExpr(e)
}

type Literal struct {
	Value interface{}
}

func (e *Literal) Accept(v Visitor) interface{} {
	return v.VisitLiteralExpr(e)
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

func (e *Unary) Accept(v Visitor) interface{} {
	return v.VisitUnaryExpr(e)
}
