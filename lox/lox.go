package lox

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type syntaxError struct {
	line    int
	message string
}

func (e syntaxError) Error() string {
	return fmt.Sprintf("syntax error on line %d: %s", e.line, e.message)
}

var hadError bool = false

func RunFile(path string) {
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
}

func RunPrompt() {
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
	p := NewParser(NewScanner(reader))
	if expr, err := p.Parse(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		hadError = true
	} else {
		fmt.Println(expr.Accept(AstPrinter).(string))
	}
}
