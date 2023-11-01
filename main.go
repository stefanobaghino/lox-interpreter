package main

import (
	"fmt"
	"lox/lox"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: lox [`demo-ast-printer` | path/to/script.lox]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		if os.Args[1] == "demo-ast-printer" {
			expression := &lox.Binary{
				Left: &lox.Unary{
					Operator: lox.Token{Type: lox.TokenType(lox.MINUS), Lexeme: "-", Literal: nil, Line: 1},
					Right: &lox.Literal{
						Value: 123,
					},
				},
				Operator: lox.Token{Type: lox.TokenType(lox.STAR), Lexeme: "*", Literal: nil, Line: 1},
				Right: &lox.Grouping{
					Expression: &lox.Literal{
						Value: 45.67,
					},
				},
			}
			fmt.Println(expression.Accept(lox.AstPrinter))
		} else {
			lox.RunFile(os.Args[1])
		}
	} else {
		lox.RunPrompt()
	}
}
