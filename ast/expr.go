package ast

import "lox/token"

type ExprVisitor interface {
	VisitBinaryExpr(*BinaryExpr) interface{}
	VisitGroupingExpr(*GroupingExpr) interface{}
	VisitLiteralExpr(*LiteralExpr) interface{}
	VisitUnaryExpr(*UnaryExpr) interface{}
}

type Expr interface {
	Accept(ExprVisitor) interface{}
}

type BinaryExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (e *BinaryExpr) Accept(v ExprVisitor) interface{} {
	return v.VisitBinaryExpr(e)
}

type GroupingExpr struct {
	Expression Expr
}

func (e *GroupingExpr) Accept(v ExprVisitor) interface{} {
	return v.VisitGroupingExpr(e)
}

type LiteralExpr struct {
	Value interface{}
}

func (e *LiteralExpr) Accept(v ExprVisitor) interface{} {
	return v.VisitLiteralExpr(e)
}

type UnaryExpr struct {
	Operator token.Token
	Right    Expr
}

func (e *UnaryExpr) Accept(v ExprVisitor) interface{} {
	return v.VisitUnaryExpr(e)
}
