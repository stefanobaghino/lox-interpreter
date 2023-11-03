package ast

import (
	"fmt"
	"strings"
)

type printer struct{}

func Print(expr Expr) string {
	return expr.Accept(printer{}).(string)
}

func (printer) VisitBinaryExpr(expr *BinaryExpr) interface{} {
	return parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}
func (printer) VisitGroupingExpr(expr *GroupingExpr) interface{} {
	return parenthesize("group", expr.Expression)
}
func (printer) VisitLiteralExpr(expr *LiteralExpr) interface{} {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}
func (printer) VisitUnaryExpr(expr *UnaryExpr) interface{} {
	return parenthesize(expr.Operator.Lexeme, expr.Right)
}

func parenthesize(name string, exprs ...Expr) string {
	builder := strings.Builder{}
	builder.WriteRune('(')
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteRune(' ')
		builder.WriteString(Print(expr))
	}
	builder.WriteRune(')')
	return builder.String()
}
