package lox

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var hadError bool = false

func RunFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	run(string(bytes))
	if hadError {
		os.Exit(65)
	}
}

func RunPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		if len(line) == 0 {
			break
		}
		run(line)
		hadError = false
	}
}

func run(source string) {
	s := NewScanner([]byte(source))
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