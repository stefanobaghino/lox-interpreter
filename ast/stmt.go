package ast

type StmtVisitor interface {
	VisitExprStmt(*ExprStmt) interface{}
	VisitPrintStmt(*PrintStmt) interface{}
	VisitEndStmt(*EndStmt) interface{}
}

type Stmt interface {
	AcceptStmt(StmtVisitor) interface{}
}

type ExprStmt struct {
	Expression Expr
}

func (s *ExprStmt) AcceptStmt(v StmtVisitor) interface{} {
	return v.VisitExprStmt(s)
}

type PrintStmt struct {
	Expression Expr
}

func (s *PrintStmt) AcceptStmt(v StmtVisitor) interface{} {
	return v.VisitPrintStmt(s)
}

type EndStmt struct {
}

func (s *EndStmt) AcceptStmt(v StmtVisitor) interface{} {
	return v.VisitEndStmt(s)
}
