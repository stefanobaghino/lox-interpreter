package main

import (
	"bufio"
	"fmt"
	"lox/runner"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: lox [path/to/script.lox]")
		os.Exit(64)
	}
	in := os.Stdin
	var exec runner.Mode = &runner.Repl{}
	if len(os.Args) == 2 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		defer file.Close()
		in = file
		exec = &runner.Script{}
	}
	runner.Run(bufio.NewReader(in), exec)
}
