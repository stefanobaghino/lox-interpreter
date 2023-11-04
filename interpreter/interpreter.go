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
	env  *Env
	done bool
}

func NewInterpreter() *Interpreter {
	return &Interpreter{env: NewGlobalEnv()}
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

func (i *Interpreter) VisitVarDeclStmt(stmt *ast.VarDeclStmt) interface{} {
	i.env.Define(stmt.Name.Lexeme, func() interface{} {
		if *stmt.Initializer == nil {
			return nil
		}
		initializer := *stmt.Initializer
		return initializer.AcceptExpr(i)
	})
	return nil
}

func (i *Interpreter) VisitBlockStmt(stmt *ast.BlockStmt) interface{} {
	env := NewEnv(i.env)
	i.env = env
	for _, stmt := range stmt.Statements {
		stmt.AcceptStmt(i)
	}
	i.env = env.parent
	return nil
}

func (i *Interpreter) VisitIfStmt(stmt *ast.IfStmt) interface{} {
	if truthy(stmt.Condition.AcceptExpr(i)) {
		(*stmt.ThenBranch).AcceptStmt(i)
	} else if stmt.ElseBranch != nil {
		(*stmt.ElseBranch).AcceptStmt(i)
	}
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *ast.PrintStmt) interface{} {
	fmt.Println(stmt.Expression.AcceptExpr(i))
	return nil
}

func (i *Interpreter) VisitAssertStmt(stmt *ast.AssertStmt) interface{} {
	assertion := stmt.Expression.AcceptExpr(i)
	if !truthy(assertion) {
		panic(&RuntimeError{line: 0 /*TODO*/, message: "assertion failed"})
	}
	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt *ast.WhileStmt) interface{} {
	for truthy(stmt.Condition.AcceptExpr(i)) {
		stmt.Body.AcceptStmt(i)
	}
	return nil
}

func (i *Interpreter) VisitEndStmt(stmt *ast.EndStmt) interface{} {
	i.done = true
	return nil
}

func (i *Interpreter) VisitExprStmt(stmt *ast.ExprStmt) interface{} {
	return stmt.Expression.AcceptExpr(i)
}

func (i *Interpreter) VisitAssignmentExpr(expr *ast.AssignmentExpr) interface{} {
	var value interface{}
	i.env.Assign(expr.Name.Lexeme, func() interface{} {
		value = expr.Value.AcceptExpr(i)
		return value
	})
	return value
}

func (i *Interpreter) VisitLogicalExpr(expr *ast.LogicalExpr) interface{} {
	left := expr.Left.AcceptExpr(i)
	if expr.Operator.Type == token.OR {
		if truthy(left) {
			return left
		}
	} else {
		if !truthy(left) {
			return left
		}
	}
	return expr.Right.AcceptExpr(i)
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

func truthy(value interface{}) bool {
	switch x := value.(type) {
	case nil:
		return false
	case bool:
		return x
	default:
		return true
	}
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
		return !truthy(right)
	}
	panic(fmt.Errorf("unexpected operator: %v", expr.Operator))
}

func (i *Interpreter) VisitVarExpr(expr *ast.VarExpr) interface{} {
	return i.env.Get(expr.Name.Lexeme)
}
