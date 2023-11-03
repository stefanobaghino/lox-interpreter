package ast

import "lox/token"

type ExprVisitor interface {
	VisitBinaryExpr(*BinaryExpr) interface{}
	VisitGroupingExpr(*GroupingExpr) interface{}
	VisitLiteralExpr(*LiteralExpr) interface{}
	VisitUnaryExpr(*UnaryExpr) interface{}
}

type Expr interface {
	AcceptExpr(ExprVisitor) interface{}
}

type BinaryExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (e *BinaryExpr) AcceptExpr(v ExprVisitor) interface{} {
	return v.VisitBinaryExpr(e)
}

type GroupingExpr struct {
	Expression Expr
}

func (e *GroupingExpr) AcceptExpr(v ExprVisitor) interface{} {
	return v.VisitGroupingExpr(e)
}

type LiteralExpr struct {
	Value interface{}
}

func (e *LiteralExpr) AcceptExpr(v ExprVisitor) interface{} {
	return v.VisitLiteralExpr(e)
}

type UnaryExpr struct {
	Operator token.Token
	Right    Expr
}

func (e *UnaryExpr) AcceptExpr(v ExprVisitor) interface{} {
	return v.VisitUnaryExpr(e)
}
