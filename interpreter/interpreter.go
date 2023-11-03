package interpreter

import (
	"fmt"
	"lox/ast"
	"lox/token"
)

type RuntimeError struct {
	line    int
	message string
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("runtime error on line %d: %s", e.line, e.message)
}

type Interpreter struct {
	done bool
}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) Interpret(stmt ast.Stmt) (result interface{}, err error) {

	defer func() {
		if e := recover(); e != nil {
			if re, ok := e.(*RuntimeError); ok {
				err = re
			} else {
				panic(fmt.Errorf("unexpected error during interpretation: %v", e))
			}
		}
	}()

	result = stmt.AcceptStmt(i)

	return
}

func (i *Interpreter) Done() bool {
	return i.done
}

func (i *Interpreter) VisitPrintStmt(stmt *ast.PrintStmt) interface{} {
	fmt.Println(stmt.Expression.AcceptExpr(i))
	return nil
}

func (i *Interpreter) VisitEndStmt(stmt *ast.EndStmt) interface{} {
	i.done = true
	return nil
}

func (i *Interpreter) VisitExprStmt(stmt *ast.ExprStmt) interface{} {
	return stmt.Expression.AcceptExpr(i)
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.BinaryExpr) interface{} {
	left := expr.Left.AcceptExpr(i)
	right := expr.Right.AcceptExpr(i)
	switch expr.Operator.Type {
	case token.MINUS:
		if left, ok := left.(float64); ok {
			if right, ok := right.(float64); ok {
				return left - right
			} else {
				panic(&RuntimeError{line: expr.Operator.Line, message: "right operand must be a number"})
			}
		} else {
			panic(&RuntimeError{line: expr.Operator.Line, message: "left operand must be a number"})
		}
	case token.SLASH:
		if left, ok := left.(float64); ok {
			if right, ok := right.(float64); ok {
				return left / right
			} else {
				panic(&RuntimeError{line: expr.Operator.Line, message: "right operand must be a number"})
			}
		} else {
			panic(&RuntimeError{line: expr.Operator.Line, message: "left operand must be a number"})
		}
	case token.STAR:
		if left, ok := left.(float64); ok {
			if right, ok := right.(float64); ok {
				return left * right
			} else {
				panic(&RuntimeError{line: expr.Operator.Line, message: "right operand must be a number"})
			}
		} else {
			panic(&RuntimeError{line: expr.Operator.Line, message: "left operand must be a number"})
		}
	case token.PLUS:
		if left, ok := left.(float64); ok {
			if right, ok := right.(float64); ok {
				return left + right
			} else {
				panic(&RuntimeError{line: expr.Operator.Line, message: "right operand must be a number"})
			}
		}
		if left, ok := left.(string); ok {
			if right, ok := right.(string); ok {
				return left + right
			} else {
				panic(&RuntimeError{line: expr.Operator.Line, message: "right operand must be a string"})
			}
		}
		panic(&RuntimeError{line: expr.Operator.Line, message: "left operand must be a number or a string"})
	case token.BANG_EQUAL:
		if left == nil {
			return right != nil
		}
		return left != right
	case token.EQUAL_EQUAL:
		if left == nil {
			return right == nil
		}
		return left == right
	case token.GREATER:
		if left, ok := left.(float64); ok {
			if right, ok := right.(float64); ok {
				return left > right
			} else {
				panic(&RuntimeError{line: expr.Operator.Line, message: "right operand must be a number"})
			}
		} else {
			panic(&RuntimeError{line: expr.Operator.Line, message: "left operand must be a number"})
		}
	case token.GREATER_EQUAL:
		if left, ok := left.(float64); ok {
			if right, ok := right.(float64); ok {
				return left >= right
			} else {
				panic(&RuntimeError{line: expr.Operator.Line, message: "right operand must be a number"})
			}
		} else {
			panic(&RuntimeError{line: expr.Operator.Line, message: "left operand must be a number"})
		}
	case token.LESS:
		if left, ok := left.(float64); ok {
			if right, ok := right.(float64); ok {
				return left < right
			} else {
				panic(&RuntimeError{line: expr.Operator.Line, message: "right operand must be a number"})
			}
		} else {
			panic(&RuntimeError{line: expr.Operator.Line, message: "left operand must be a number"})
		}
	case token.LESS_EQUAL:
		if left, ok := left.(float64); ok {
			if right, ok := right.(float64); ok {
				return left <= right
			} else {
				panic(&RuntimeError{line: expr.Operator.Line, message: "right operand must be a number"})
			}
		} else {
			panic(&RuntimeError{line: expr.Operator.Line, message: "left operand must be a number"})
		}
	}
	panic(fmt.Errorf("unexpected operator: %v", expr.Operator))
}

func (i *Interpreter) VisitGroupingExpr(expr *ast.GroupingExpr) interface{} {
	return expr.Expression.AcceptExpr(i)
}

func (i *Interpreter) VisitLiteralExpr(expr *ast.LiteralExpr) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.UnaryExpr) interface{} {
	right := expr.Right.AcceptExpr(i)
	switch expr.Operator.Type {
	case token.MINUS:
		if x, ok := right.(float64); ok {
			return -x
		} else {
			panic(&RuntimeError{line: expr.Operator.Line, message: "operand must be a number"})
		}
	case token.BANG:
		switch x := right.(type) {
		case nil:
			return true
		case bool:
			return !x
		default:
			return false
		}
	}
	panic(fmt.Errorf("unexpected operator: %v", expr.Operator))
}
