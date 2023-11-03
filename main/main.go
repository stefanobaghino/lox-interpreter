package main

import (
	"bufio"
	"fmt"
	"lox/interpreter"
	"lox/parser"
	"lox/scanner"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: lox [path/to/script.lox]")
		os.Exit(64)
	}
	in := os.Stdin
	var exec executionMode = &replMode{}
	if len(os.Args) == 2 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		defer file.Close()
		in = file
		exec = &scriptMode{}
	}
	run(bufio.NewReader(in), exec)
}

func run(reader *bufio.Reader, exec executionMode) {
	p := parser.NewParser(scanner.NewScanner(reader))
	i := interpreter.NewInterpreter()
	for {
		exec.PreStmt()
		if stmt, err := p.NextStatement(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			exec.PostError(err)
		} else {
			res, err := i.Interpret(stmt)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				exec.PostError(err)
			} else if !i.Done() {
				exec.PostStmt(res)
			} else {
				break
			}
		}
	}
}
