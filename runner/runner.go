package runner

import (
	"bufio"
	"fmt"
	"lox/interpreter"
	"lox/parser"
	"lox/resolver"
	"lox/scanner"
	"os"
)

func Run(reader *bufio.Reader, mode Mode) {
	p := parser.NewParser(scanner.NewScanner(reader))
	i := interpreter.NewInterpreter()
	r := resolver.NewResolver(i)
	for {
		mode.PreStmt()
		if stmt, err := p.NextStatement(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			mode.PostGrammarError(err)
		} else if mode.Execute() {
			err := r.Resolve(stmt)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				mode.PostGrammarError(err)
			}
			res, err := i.Interpret(stmt)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				mode.PostRuntimeError(err)
			} else if i.Done() {
				break
			}
			mode.PostStmt(res)
		} else {
			break
		}
	}
}
