package lox

import "fmt"

type runtimeError struct {
	line    int
	message string
}

func (e runtimeError) Error() string {
	return fmt.Sprintf("runtime error on line %d: %s", e.line, e.message)
}

type Interpreter struct{}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) Interpret(expr Expr) (result interface{}, err error) {

	defer func() {
		if e := recover(); e != nil {
			if re, ok := e.(*runtimeError); ok {
				err = re
			} else {
				panic(fmt.Errorf("unexpected error during interpretation: %v", e))
			}
		}
	}()

	result = expr.Accept(i)

	return
}

func (i *Interpreter) VisitBinaryExpr(expr *Binary) interface{} {
	left := expr.Left.Accept(i)
	right := expr.Right.Accept(i)
	switch expr.Operator.Type {
	case MINUS:
		if left, ok := left.(float64); ok {
			if right, ok := right.(float64); ok {
				return left - right
			} else {
				panic(&runtimeError{line: expr.Operator.Line, message: "right operand must be a number"})
			}
		} else {
			panic(&runtimeError{line: expr.Operator.Line, message: "left operand must be a number"})
		}
	case SLASH:
		if left, ok := left.(float64); ok {
			if right, ok := right.(float64); ok {
				return left / right
			} else {
				panic(&runtimeError{line: expr.Operator.Line, message: "right operand must be a number"})
			}
		} else {
			panic(&runtimeError{line: expr.Operator.Line, message: "left operand must be a number"})
		}
	case STAR:
		if left, ok := left.(float64); ok {
			if right, ok := right.(float64); ok {
				return left * right
			} else {
				panic(&runtimeError{line: expr.Operator.Line, message: "right operand must be a number"})
			}
		} else {
			panic(&runtimeError{line: expr.Operator.Line, message: "left operand must be a number"})
		}
	case PLUS:
		if left, ok := left.(float64); ok {
			if right, ok := right.(float64); ok {
				return left + right
			} else {
				panic(&runtimeError{line: expr.Operator.Line, message: "right operand must be a number"})
			}
		}
		if left, ok := left.(string); ok {
			if right, ok := right.(string); ok {
				return left + right
			} else {
				panic(&runtimeError{line: expr.Operator.Line, message: "right operand must be a string"})
			}
		}
		panic(&runtimeError{line: expr.Operator.Line, message: "left operand must be a number or a string"})
	case BANG_EQUAL:
		if left == nil {
			return right != nil
		}
		return left != right
	case EQUAL_EQUAL:
		if left == nil {
			return right == nil
		}
		return left == right
	case GREATER:
		if left, ok := left.(float64); ok {
			if right, ok := right.(float64); ok {
				return left > right
			} else {
				panic(&runtimeError{line: expr.Operator.Line, message: "right operand must be a number"})
			}
		} else {
			panic(&runtimeError{line: expr.Operator.Line, message: "left operand must be a number"})
		}
	case GREATER_EQUAL:
		if left, ok := left.(float64); ok {
			if right, ok := right.(float64); ok {
				return left >= right
			} else {
				panic(&runtimeError{line: expr.Operator.Line, message: "right operand must be a number"})
			}
		} else {
			panic(&runtimeError{line: expr.Operator.Line, message: "left operand must be a number"})
		}
	case LESS:
		if left, ok := left.(float64); ok {
			if right, ok := right.(float64); ok {
				return left < right
			} else {
				panic(&runtimeError{line: expr.Operator.Line, message: "right operand must be a number"})
			}
		} else {
			panic(&runtimeError{line: expr.Operator.Line, message: "left operand must be a number"})
		}
	case LESS_EQUAL:
		if left, ok := left.(float64); ok {
			if right, ok := right.(float64); ok {
				return left <= right
			} else {
				panic(&runtimeError{line: expr.Operator.Line, message: "right operand must be a number"})
			}
		} else {
			panic(&runtimeError{line: expr.Operator.Line, message: "left operand must be a number"})
		}
	}
	panic(fmt.Errorf("unexpected operator: %v", expr.Operator))
}

func (i *Interpreter) VisitGroupingExpr(expr *Grouping) interface{} {
	return expr.Expression.Accept(i)
}

func (i *Interpreter) VisitLiteralExpr(expr *Literal) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitUnaryExpr(expr *Unary) interface{} {
	right := expr.Right.Accept(i)
	switch expr.Operator.Type {
	case MINUS:
		if x, ok := right.(float64); ok {
			return -x
		} else {
			panic(&runtimeError{line: expr.Operator.Line, message: "operand must be a number"})
		}
	case BANG:
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
