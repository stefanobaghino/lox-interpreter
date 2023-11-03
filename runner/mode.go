package runner

import (
	"fmt"
	"os"
)

type Mode interface {
	PreStmt()
	PostStmt(interface{})
	PostError(error)
}

type Repl struct{}

func (m *Repl) PreStmt() {
	fmt.Print("> ")
}

func (m *Repl) PostStmt(res interface{}) {
	fmt.Printf("%v\n", res)
}

func (m *Repl) PostError(err error) {
}

type Script struct{}

func (m *Script) PreStmt() {
}

func (m *Script) PostStmt(res interface{}) {
}

func (m *Script) PostError(err error) {
	os.Exit(1)
}
