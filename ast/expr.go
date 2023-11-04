package ast

import "lox/token"

type ExprVisitor interface {
	VisitAssignmentExpr(*AssignmentExpr) interface{}
	VisitBinaryExpr(*BinaryExpr) interface{}
	VisitGroupingExpr(*GroupingExpr) interface{}
	VisitLiteralExpr(*LiteralExpr) interface{}
	VisitLogicalExpr(*LogicalExpr) interface{}
	VisitUnaryExpr(*UnaryExpr) interface{}
	VisitVarExpr(*VarExpr) interface{}
}

type Expr interface {
	AcceptExpr(ExprVisitor) interface{}
}

type AssignmentExpr struct {
	Name  token.Token
	Value Expr
}

func (e *AssignmentExpr) AcceptExpr(v ExprVisitor) interface{} {
	return v.VisitAssignmentExpr(e)
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

type LogicalExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (e *LogicalExpr) AcceptExpr(v ExprVisitor) interface{} {
	return v.VisitLogicalExpr(e)
}

type UnaryExpr struct {
	Operator token.Token
	Right    Expr
}

func (e *UnaryExpr) AcceptExpr(v ExprVisitor) interface{} {
	return v.VisitUnaryExpr(e)
}

type VarExpr struct {
	Name token.Token
}

func (e *VarExpr) AcceptExpr(v ExprVisitor) interface{} {
	return v.VisitVarExpr(e)
}
