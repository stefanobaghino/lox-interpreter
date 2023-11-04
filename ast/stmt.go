package ast

import "lox/token"

type StmtVisitor interface {
	VisitVarDeclStmt(*VarDeclStmt) interface{}
	VisitFunDeclStmt(*FunDeclStmt) interface{}
	VisitBlockStmt(*BlockStmt) interface{}
	VisitExprStmt(*ExprStmt) interface{}
	VisitIfStmt(*IfStmt) interface{}
	VisitAssertStmt(*AssertStmt) interface{}
	VisitPrintStmt(*PrintStmt) interface{}
	VisitWhileStmt(*WhileStmt) interface{}
	VisitEndStmt(*EndStmt) interface{}
}

type Stmt interface {
	AcceptStmt(StmtVisitor) interface{}
}

type VarDeclStmt struct {
	Name        token.Token
	Initializer *Expr
}

func (s *VarDeclStmt) AcceptStmt(v StmtVisitor) interface{} {
	return v.VisitVarDeclStmt(s)
}

type FunDeclStmt struct {
	Name   token.Token
	Params []token.Token
	Body   *BlockStmt
}

func (s *FunDeclStmt) AcceptStmt(v StmtVisitor) interface{} {
	return v.VisitFunDeclStmt(s)
}

type BlockStmt struct {
	Statements []Stmt
}

func (s *BlockStmt) AcceptStmt(v StmtVisitor) interface{} {
	return v.VisitBlockStmt(s)
}

type ExprStmt struct {
	Expression Expr
}

func (s *ExprStmt) AcceptStmt(v StmtVisitor) interface{} {
	return v.VisitExprStmt(s)
}

type IfStmt struct {
	Condition  Expr
	ThenBranch *Stmt
	ElseBranch *Stmt
}

func (s *IfStmt) AcceptStmt(v StmtVisitor) interface{} {
	return v.VisitIfStmt(s)
}

type AssertStmt struct {
	Expression Expr
}

func (s *AssertStmt) AcceptStmt(v StmtVisitor) interface{} {
	return v.VisitAssertStmt(s)
}

type PrintStmt struct {
	Expression Expr
}

func (s *PrintStmt) AcceptStmt(v StmtVisitor) interface{} {
	return v.VisitPrintStmt(s)
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (s *WhileStmt) AcceptStmt(v StmtVisitor) interface{} {
	return v.VisitWhileStmt(s)
}

type EndStmt struct {
}

func (s *EndStmt) AcceptStmt(v StmtVisitor) interface{} {
	return v.VisitEndStmt(s)
}
