package lox

import (
	"fmt"
	"strings"
)

type astPrinter struct{}

var AstPrinter = astPrinter{}

func (astPrinter) VisitBinaryExpr(expr *Binary) interface{} {
	return parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}
func (astPrinter) VisitGroupingExpr(expr *Grouping) interface{} {
	return parenthesize("group", expr.Expression)
}
func (astPrinter) VisitLiteralExpr(expr *Literal) interface{} {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}
func (astPrinter) VisitUnaryExpr(expr *Unary) interface{} {
	return parenthesize(expr.Operator.Lexeme, expr.Right)
}

func parenthesize(name string, exprs ...Expr) string {
	builder := strings.Builder{}
	builder.WriteRune('(')
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteRune(' ')
		builder.WriteString(expr.Accept(AstPrinter).(string))
	}
	builder.WriteRune(')')
	return builder.String()
}
