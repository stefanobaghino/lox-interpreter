package main

import (
	"fmt"
	"os"
)

type executionMode interface {
	PreStmt()
	PostStmt(interface{})
	PostError(error)
}

type replMode struct{}

func (m *replMode) PreStmt() {
	fmt.Print("> ")
}

func (m *replMode) PostStmt(res interface{}) {
	fmt.Printf("%v\n", res)
}

func (m *replMode) PostError(err error) {
}

type scriptMode struct{}

func (m *scriptMode) PreStmt() {
}

func (m *scriptMode) PostStmt(res interface{}) {
}

func (m *scriptMode) PostError(err error) {
	os.Exit(1)
}
