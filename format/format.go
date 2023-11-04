package format

import (
	"fmt"
	"lox/ast"
	"strings"
)

type Formatter struct {
	indent int
}

func (f *Formatter) fmtExpr(expr ast.Expr) string {
	return expr.AcceptExpr(f).(string)
}

func NewFormatter() *Formatter {
	return &Formatter{}
}

func (f *Formatter) Format(stmt ast.Stmt) string {
	builder := strings.Builder{}
	for i := 0; i < f.indent; i++ {
		builder.WriteString("\t")
	}
	builder.WriteString(stmt.AcceptStmt(f).(string))
	return builder.String()
}

func (f *Formatter) VisitVarDeclStmt(stmt *ast.VarDeclStmt) interface{} {
	builder := strings.Builder{}
	builder.WriteString("var ")
	builder.WriteString(stmt.Name.Lexeme)
	if stmt.Initializer != nil {
		builder.WriteString(" = ")
		builder.WriteString(f.fmtExpr(*stmt.Initializer))
	}
	builder.WriteRune(';')
	return builder.String()
}

func (f *Formatter) VisitBlockStmt(stmt *ast.BlockStmt) interface{} {
	builder := strings.Builder{}
	builder.WriteRune('\n')
	for i := 0; i < f.indent; i++ {
		builder.WriteString("\t")
	}
	builder.WriteRune('{')
	f.indent++
	for _, stmt := range stmt.Statements {
		builder.WriteRune('\n')
		builder.WriteString(f.Format(stmt))
	}
	builder.WriteRune('\n')
	builder.WriteRune('}')
	return builder.String()
}

func (f *Formatter) VisitExprStmt(stmt *ast.ExprStmt) interface{} {
	builder := strings.Builder{}
	builder.WriteString(f.fmtExpr(stmt.Expression))
	builder.WriteRune(';')
	return builder.String()
}

func (f *Formatter) VisitPrintStmt(stmt *ast.PrintStmt) interface{} {
	builder := strings.Builder{}
	builder.WriteString("print ")
	builder.WriteString(f.fmtExpr(stmt.Expression))
	builder.WriteRune(';')
	return builder.String()
}

func (f *Formatter) VisitAssertStmt(stmt *ast.AssertStmt) interface{} {
	builder := strings.Builder{}
	builder.WriteString("assert ")
	builder.WriteString(f.fmtExpr(stmt.Expression))
	builder.WriteRune(';')
	return builder.String()
}

func (f *Formatter) VisitEndStmt(stmt *ast.EndStmt) interface{} {
	return ""
}

func (f *Formatter) VisitAssignmentExpr(expr *ast.AssignmentExpr) interface{} {
	builder := strings.Builder{}
	builder.WriteString(expr.Name.Lexeme)
	builder.WriteString(" = ")
	builder.WriteString(f.fmtExpr(expr.Value))
	return builder.String()
}

func (f *Formatter) VisitBinaryExpr(expr *ast.BinaryExpr) interface{} {
	builder := strings.Builder{}
	builder.WriteString(f.fmtExpr(expr.Left))
	builder.WriteRune(' ')
	builder.WriteString(expr.Operator.Lexeme)
	builder.WriteRune(' ')
	builder.WriteString(f.fmtExpr(expr.Right))
	return builder.String()
}

func (f *Formatter) VisitGroupingExpr(expr *ast.GroupingExpr) interface{} {
	builder := strings.Builder{}
	builder.WriteRune('(')
	builder.WriteString(f.fmtExpr(expr.Expression))
	builder.WriteRune(')')
	return builder.String()
}

func (f *Formatter) VisitLiteralExpr(expr *ast.LiteralExpr) interface{} {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (f *Formatter) VisitUnaryExpr(expr *ast.UnaryExpr) interface{} {
	builder := strings.Builder{}
	builder.WriteString(expr.Operator.Lexeme)
	builder.WriteString(f.fmtExpr(expr.Right))
	return builder.String()
}

func (f *Formatter) VisitVarExpr(expr *ast.VarExpr) interface{} {
	return expr.Name.Lexeme
}
