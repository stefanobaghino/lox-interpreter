package format

import (
	"fmt"
	"lox/ast"
	"strings"
)

type Formatter struct {
	indentation int
}

func (f *Formatter) fmtExpr(expr ast.Expr) string {
	return expr.AcceptExpr(f).(string)
}

func NewFormatter() *Formatter {
	return &Formatter{}
}

func (f *Formatter) Format(stmt ast.Stmt) string {
	builder := strings.Builder{}
	f.indent(&builder)
	builder.WriteString(stmt.AcceptStmt(f).(string))
	return builder.String()
}

func (f *Formatter) block(stmts func() []ast.Stmt) string {
	builder := strings.Builder{}
	builder.WriteRune('\n')
	f.indent(&builder)
	builder.WriteRune('{')
	f.indentation++
	for _, stmt := range stmts() {
		builder.WriteRune('\n')
		builder.WriteString(f.Format(stmt))
	}
	builder.WriteRune('\n')
	builder.WriteRune('}')
	return builder.String()
}

func (f *Formatter) indent(builder *strings.Builder) {
	for i := 0; i < f.indentation; i++ {
		builder.WriteRune('\t')
	}
}

func (f *Formatter) VisitFunDeclStmt(stmt *ast.FunDeclStmt) interface{} {
	builder := strings.Builder{}
	builder.WriteString("fun ")
	builder.WriteString(stmt.Name.Lexeme)
	builder.WriteRune('(')
	for i, param := range stmt.Params {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(param.Lexeme)
	}
	builder.WriteRune(')')
	builder.WriteString(f.Format(stmt.Body))
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
	return f.block(func() []ast.Stmt { return stmt.Statements })
}

func (f *Formatter) VisitExprStmt(stmt *ast.ExprStmt) interface{} {
	builder := strings.Builder{}
	builder.WriteString(f.fmtExpr(stmt.Expression))
	builder.WriteRune(';')
	return builder.String()
}

func (f *Formatter) VisitIfStmt(stmt *ast.IfStmt) interface{} {
	builder := strings.Builder{}
	builder.WriteString("if (")
	builder.WriteString(f.fmtExpr(stmt.Condition))
	builder.WriteRune(')')
	thenBlock := *stmt.ThenBranch
	if _, ok := thenBlock.(*ast.BlockStmt); !ok {
		thenBlock = &ast.BlockStmt{Statements: []ast.Stmt{thenBlock}}
	}
	builder.WriteString(f.Format(thenBlock))
	if stmt.ElseBranch != nil {
		builder.WriteRune('\n')
		f.indent(&builder)
		builder.WriteString("else")

		elseBlock := *stmt.ElseBranch
		if _, ok := elseBlock.(*ast.BlockStmt); !ok {
			elseBlock = &ast.BlockStmt{Statements: []ast.Stmt{elseBlock}}
		}
		builder.WriteString(f.Format(elseBlock))
	}
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

func (f *Formatter) VisitWhileStmt(stmt *ast.WhileStmt) interface{} {
	builder := strings.Builder{}
	builder.WriteString("while (")
	builder.WriteString(f.fmtExpr(stmt.Condition))
	builder.WriteRune(')')
	builder.WriteString(f.Format(stmt.Body))
	return builder.String()
}

func (f *Formatter) VisitAssignmentExpr(expr *ast.AssignmentExpr) interface{} {
	builder := strings.Builder{}
	builder.WriteString(expr.Name.Lexeme)
	builder.WriteString(" = ")
	builder.WriteString(f.fmtExpr(expr.Value))
	return builder.String()
}

func (f *Formatter) VisitLogicalExpr(expr *ast.LogicalExpr) interface{} {
	builder := strings.Builder{}
	builder.WriteString(f.fmtExpr(expr.Left))
	builder.WriteRune(' ')
	builder.WriteString(expr.Operator.Lexeme)
	builder.WriteRune(' ')
	builder.WriteString(f.fmtExpr(expr.Right))
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

func (f *Formatter) VisitCallExpr(expr *ast.CallExpr) interface{} {
	builder := strings.Builder{}
	builder.WriteString(f.fmtExpr(expr.Callee))
	builder.WriteRune('(')
	for i, arg := range expr.Arguments {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(f.fmtExpr(arg))
	}
	builder.WriteRune(')')
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
	switch expr.Value.(type) {
	case string:
		return fmt.Sprintf("%q", expr.Value)
	case float64:
		return fmt.Sprintf("%f", expr.Value)
	case nil:
		return "nil"
	default:
		return fmt.Sprintf("%v", expr.Value)
	}
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
