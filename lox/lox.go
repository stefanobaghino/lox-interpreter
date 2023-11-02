package lox

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

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
	s := NewScanner(reader)
	for {
		t, e := s.NextToken()
		if e != nil {
			fmt.Fprintln(os.Stderr, e.Error())
			hadError = true
		}
		if t.Type == EOF {
			break
		}
		fmt.Println(t)
	}
}
