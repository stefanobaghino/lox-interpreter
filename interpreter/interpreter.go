package interpreter

import (
	"fmt"
	"lox/ast"
	"lox/token"
	"time"
)

type Callable interface {
	Arity() int
	Call(*Interpreter, []interface{}) interface{}
}

type function struct {
	arity   int
	closure *Env
	call    func(*Interpreter, []interface{}) interface{}
}

func newFunction(arity int, closure *Env, call func(*Interpreter, []interface{}) interface{}) *function {
	return &function{arity: arity, closure: closure, call: call}
}

func (b *function) Arity() int {
	return b.arity
}

func (b *function) Call(i *Interpreter, arguments []interface{}) interface{} {
	i.env = b.closure
	result := b.call(i, arguments)
	return result
}

type Return struct {
	value interface{}
}

type RuntimeError struct {
	line    int
	message string
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("runtime error on line %d: %s", e.line, e.message)
}

type Interpreter struct {
	locals  map[ast.Expr]int
	globals *Env
	env     *Env
	done    bool
}

func NewInterpreter() *Interpreter {
	globals := NewGlobalEnv()
	globals.Define("clock", func() interface{} {
		return newFunction(0, globals, func(i *Interpreter, arguments []interface{}) interface{} {
			return float64(time.Now().Unix())
		})
	})
	return &Interpreter{locals: make(map[ast.Expr]int), globals: globals, env: globals}
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

func (i *Interpreter) Resolve(expr ast.Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) VisitFunDeclStmt(stmt *ast.FunDeclStmt) interface{} {
	i.env.Define(stmt.Name.Lexeme, func() interface{} {
		return newFunction(len(stmt.Params), i.env, func(i *Interpreter, arguments []interface{}) (ret interface{}) {
			env := NewEnv(i.env)
			i.env = env
			for index, param := range stmt.Params {
				env.Define(param.Lexeme, func() interface{} {
					return arguments[index]
				})
			}
			defer func() {
				i.env = env.parent
				if e := recover(); e != nil {
					if r, ok := e.(*Return); ok {
						ret = r.value
					} else {
						panic(e)
					}
				}
			}()
			stmt.Body.AcceptStmt(i)
			return
		})
	})
	return nil
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

func (i *Interpreter) VisitReturnStmt(stmt *ast.ReturnStmt) interface{} {
	var value interface{}
	if stmt.Value != nil {
		value = (*stmt.Value).AcceptExpr(i)
	}
	panic(&Return{value: value})
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
	if distance, ok := i.locals[expr]; ok {
		value = i.env.AssignAt(distance, expr.Name.Lexeme, func() interface{} {
			value = expr.Value.AcceptExpr(i)
			return value
		})
	} else {
		value = i.globals.Assign(expr.Name.Lexeme, func() interface{} {
			value = expr.Value.AcceptExpr(i)
			return value
		})
	}
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

func (i *Interpreter) VisitCallExpr(expr *ast.CallExpr) interface{} {
	callee := expr.Callee.AcceptExpr(i)
	var arguments []interface{}
	for _, arg := range expr.Arguments {
		arguments = append(arguments, arg.AcceptExpr(i))
	}
	if function, ok := callee.(Callable); ok {
		if len(arguments) != function.Arity() {
			panic(&RuntimeError{line: expr.Paren.Line, message: fmt.Sprintf("expected %d arguments but got %d", function.Arity(), len(arguments))})
		}
		return function.Call(i, arguments)
	} else {
		panic(&RuntimeError{line: expr.Paren.Line, message: "identifier is not a function"})
	}
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
	return i.lookupVariable(expr.Name, expr)
}

func (i *Interpreter) lookupVariable(name token.Token, expr ast.Expr) interface{} {
	if distance, ok := i.locals[expr]; ok {
		return i.env.GetAt(distance, name.Lexeme)
	} else {
		return i.globals.Get(name.Lexeme)
	}
}
