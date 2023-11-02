package main

import (
	"fmt"
	"lox/lox"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: lox [path/to/script.lox]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		lox.RunFile(os.Args[1])
	} else {
		lox.RunPrompt()
	}
}
