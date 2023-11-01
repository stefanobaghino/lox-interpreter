package main

import (
	"fmt"
	"lox/lox"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Println("Usage: lox [script]")
		os.Exit(64)
	} else if len(os.Args) == 1 {
		lox.RunFile(os.Args[0])
	} else {
		lox.RunPrompt()
	}
}
