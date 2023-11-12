package resolver

import (
	"fmt"
	"lox/ast"
	"lox/interpreter"
	"lox/token"
)

type Resolver struct {
	interpreter *interpreter.Interpreter
	scopes      []map[string]bool
}

func NewResolver(interpreter *interpreter.Interpreter) *Resolver {
	return &Resolver{interpreter: interpreter}
}

type ResolutionError struct {
	line    int
	message string
}

func (e ResolutionError) Error() string {
	return fmt.Sprintf("resolution error on line %d: %s", e.line, e.message)
}

func (r *Resolver) Resolve(stmt ast.Stmt) (err error) {
	defer func() {
		if e := recover(); e != nil {
			if re, ok := e.(*ResolutionError); ok {
				err = re
			} else {
				panic(fmt.Errorf("unexpected error during resolution: %v", e))
			}
		}
	}()
	r.resolveStmt(stmt)
	return err
}

func (r *Resolver) VisitAssignmentExpr(expr *ast.AssignmentExpr) interface{} {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitBinaryExpr(expr *ast.BinaryExpr) interface{} {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitCallExpr(expr *ast.CallExpr) interface{} {
	r.resolveExpr(expr.Callee)
	for _, arg := range expr.Arguments {
		r.resolveExpr(arg)
	}
	return nil
}

func (r *Resolver) VisitGroupingExpr(expr *ast.GroupingExpr) interface{} {
	r.resolveExpr(expr.Expression)
	return nil
}

func (r *Resolver) VisitLiteralExpr(expr *ast.LiteralExpr) interface{} {
	return nil
}

func (r *Resolver) VisitLogicalExpr(expr *ast.LogicalExpr) interface{} {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitUnaryExpr(expr *ast.UnaryExpr) interface{} {
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) resolveLocal(expr ast.Expr, name token.Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			r.interpreter.Resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) VisitVarExpr(expr *ast.VarExpr) interface{} {
	if len(r.scopes) > 0 && !r.scopes[len(r.scopes)-1][expr.Name.Lexeme] {
		panic(&ResolutionError{line: expr.Name.Line, message: "cannot read local variable in its own initializer"})
	}
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) declare(name token.Token) {
	if len(r.scopes) == 0 {
		return
	}
	if _, ok := r.scopes[len(r.scopes)-1][name.Lexeme]; ok {
		panic(&ResolutionError{line: name.Line, message: "variable with this name already declared in this scope"})
	}
	scope := r.scopes[len(r.scopes)-1]
	scope[name.Lexeme] = false
}

func (r *Resolver) define(name token.Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := r.scopes[len(r.scopes)-1]
	scope[name.Lexeme] = true
}

func (r *Resolver) VisitVarDeclStmt(stmt *ast.VarDeclStmt) interface{} {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpr(*stmt.Initializer)
	}
	r.define(stmt.Name)
	return nil
}

func (r *Resolver) resolveFunction(stmt *ast.FunDeclStmt) {
	r.beginScope()
	for _, param := range stmt.Params {
		r.declare(param)
		r.define(param)
	}
	r.resolveStmts(stmt.Body.Statements)
	r.endScope()
}

func (r *Resolver) VisitFunDeclStmt(stmt *ast.FunDeclStmt) interface{} {
	r.declare(stmt.Name)
	r.define(stmt.Name)
	r.resolveFunction(stmt)
	return nil
}

func (r *Resolver) resolveStmt(stmt ast.Stmt) {
	stmt.AcceptStmt(r)
}

func (r *Resolver) resolveExpr(expr ast.Expr) {
	expr.AcceptExpr(r)
}

func (r *Resolver) resolveStmts(stmts []ast.Stmt) {
	for _, stmt := range stmts {
		r.resolveStmt(stmt)
	}
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) VisitBlockStmt(stmt *ast.BlockStmt) interface{} {
	r.beginScope()
	r.resolveStmts(stmt.Statements)
	r.endScope()
	return nil
}

func (r *Resolver) VisitExprStmt(stmt *ast.ExprStmt) interface{} {
	r.resolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitIfStmt(stmt *ast.IfStmt) interface{} {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(*stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStmt(*stmt.ElseBranch)
	}
	return nil
}

func (r *Resolver) VisitAssertStmt(stmt *ast.AssertStmt) interface{} {
	r.resolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitPrintStmt(stmt *ast.PrintStmt) interface{} {
	r.resolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitWhileStmt(stmt *ast.WhileStmt) interface{} {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Body)
	return nil
}

func (r *Resolver) VisitReturnStmt(stmt *ast.ReturnStmt) interface{} {
	return nil
}

func (r *Resolver) VisitEndStmt(stmt *ast.EndStmt) interface{} {
	return nil
}
