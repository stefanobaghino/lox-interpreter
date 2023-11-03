package main

import (
	"bufio"
	"fmt"
	"io"
	"lox/interpreter"
	"lox/parser"
	"lox/scanner"
	"os"
	"strings"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: lox [path/to/script.lox]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

var hadError bool = false
var hadRuntimeError bool = false

func runFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	defer file.Close()
	run(bufio.NewReader(file))
	if hadError {
		os.Exit(65)
	}
	if hadRuntimeError {
		os.Exit(70)
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		run(bufio.NewReader(strings.NewReader(line)))
		hadError = false
	}
}

func run(reader *bufio.Reader) {
	p := parser.NewParser(scanner.NewScanner(reader))
	if expr, err := p.Parse(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		hadError = true
	} else {
		res, err := interpreter.NewInterpreter().Interpret(expr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			hadRuntimeError = true
		} else {
			fmt.Printf("%v\n", res)
		}
	}
}
