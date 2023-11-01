package lox

type Expr interface {
	Accept(Visitor) interface{}
}

type Visitor interface {
	VisitBinaryExpr(*Binary) interface{}
	VisitGroupingExpr(*Grouping) interface{}
	VisitLiteralExpr(*Literal) interface{}
	VisitUnaryExpr(*Unary) interface{}
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (e *Binary) Accept(v Visitor) interface{} {
	return v.VisitBinaryExpr(e)
}

type Grouping struct {
	Expression Expr
}

func (e *Grouping) Accept(v Visitor) interface{} {
	return v.VisitGroupingExpr(e)
}

type Literal struct {
	Value interface{}
}

func (e *Literal) Accept(v Visitor) interface{} {
	return v.VisitLiteralExpr(e)
}

type Unary struct {
	Operator Token
	Right    Expr
}

func (e *Unary) Accept(v Visitor) interface{} {
	return v.VisitUnaryExpr(e)
}
