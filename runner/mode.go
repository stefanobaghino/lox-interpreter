package runner

import (
	"fmt"
	"os"
)

type Mode interface {
	PreStmt()
	PostStmt(interface{})
	PostGrammarError(error)
	PostRuntimeError(error)
	Execute() bool
}

type Repl struct{}

func (m *Repl) PreStmt() {
	fmt.Print("> ")
}

func (m *Repl) PostStmt(res interface{}) {
	fmt.Printf("%v\n", res)
}

func (m *Repl) PostGrammarError(err error) {
}

func (m *Repl) PostRuntimeError(err error) {
}

func (m *Repl) Execute() bool {
	return true
}

type Script struct {
	grammarError bool
}

func (m *Script) PreStmt() {
}

func (m *Script) PostStmt(res interface{}) {
}

func (m *Script) PostGrammarError(err error) {
	m.grammarError = true
}

func (m *Script) PostRuntimeError(err error) {
	os.Exit(1)
}

func (m *Script) Execute() bool {
	return !m.grammarError
}
