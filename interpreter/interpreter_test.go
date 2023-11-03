package interpreter

import (
	"bufio"
	"lox/parser"
	"lox/scanner"
	"regexp"
	"strings"
	"testing"
)

func TestInterpreterSimpleExpr(t *testing.T) {
	expectResult(t, "1 + 2", 3.0)
	expectResult(t, "-(1 * (2 + 3) / (4 - 5))", 5.0)
	expectResult(t, "\"foot\" + \"ball\"", "football")
}

func TestInterpreterBooleanNegation(t *testing.T) {
	expectResult(t, "!true", false)
	expectResult(t, "!false", true)
	expectResult(t, "!nil", true)
	expectResult(t, "!!true", true)
	expectResult(t, "!\"hi\"", false)
}

func TestInterpreterComparisons(t *testing.T) {
	expectResult(t, "1 < 2", true)
	expectResult(t, "1 <= 2", true)
	expectResult(t, "1 > 2", false)
	expectResult(t, "1 >= 2", false)
	expectResult(t, "1 == 2", false)
	expectResult(t, "1 != 2", true)
	expectResult(t, "1 == 1", true)
	expectResult(t, "1 != 1", false)
	expectResult(t, "nil == nil", true)
	expectResult(t, "nil != nil", false)
	expectResult(t, "nil == 0", false)
	expectResult(t, "nil != 0", true)
	expectResult(t, "true == true", true)
	expectResult(t, "true != true", false)
	expectResult(t, "true == false", false)
	expectResult(t, "true != false", true)
	expectResult(t, "false == false", true)
	expectResult(t, "false != false", false)
	expectResult(t, "false == true", false)
	expectResult(t, "false != true", true)
}

func TestInterpreterGrouping(t *testing.T) {
	src := "1 * (2 + 3)"
	expectResult(t, src, 5.0)
}

func TestInterpreterTypeErrors(t *testing.T) {
	expectRuntimeError(t, "1 + \"hi\"", "right operand must be a number")
	expectRuntimeError(t, "\"hi\" + 1", "right operand must be a string")
	expectRuntimeError(t, "1 - \"hi\"", "right operand must be a number")
	expectRuntimeError(t, "\"hi\" - 1", "left operand must be a number")
	expectRuntimeError(t, "1 * \"hi\"", "right operand must be a number")
	expectRuntimeError(t, "\"hi\" * 1", "left operand must be a number")
	expectRuntimeError(t, "1 / \"hi\"", "right operand must be a number")
	expectRuntimeError(t, "\"hi\" / 1", "left operand must be a number")
	expectRuntimeError(t, "1 < \"hi\"", "right operand must be a number")
	expectRuntimeError(t, "\"hi\" < 1", "left operand must be a number")
	expectRuntimeError(t, "1 <= \"hi\"", "right operand must be a number")
	expectRuntimeError(t, "\"hi\" <= 1", "left operand must be a number")
	expectRuntimeError(t, "1 > \"hi\"", "right operand must be a number")
	expectRuntimeError(t, "\"hi\" > 1", "left operand must be a number")
	expectRuntimeError(t, "1 >= \"hi\"", "right operand must be a number")
	expectRuntimeError(t, "\"hi\" >= 1", "left operand must be a number")
	expectRuntimeError(t, "true + 1", "left operand must be a number or a string")
	expectRuntimeError(t, "true < 1", "left operand must be a number")
	expectRuntimeError(t, "1 < true", "right operand must be a number")
	expectRuntimeError(t, "true <= 1", "left operand must be a number")
	expectRuntimeError(t, "1 <= true", "right operand must be a number")
	expectRuntimeError(t, "true > 1", "left operand must be a number")
	expectRuntimeError(t, "1 > true", "right operand must be a number")
	expectRuntimeError(t, "true >= 1", "left operand must be a number")
	expectRuntimeError(t, "1 >= true", "right operand must be a number")
	expectRuntimeError(t, "-true", "operand must be a number")
}

func expectRuntimeError(t *testing.T, src string, regex string) {
	t.Helper()
	if _, err := interpret(t, src); err == nil {
		t.Errorf("expected runtime error matching '%s', got none", regex)
	} else if re, ok := err.(*RuntimeError); ok {
		if !regexp.MustCompile(regex).MatchString(re.Error()) {
			t.Errorf("expected runtime error matching '%s', got '%v'", regex, re.Error())
		}
	} else {
		t.Errorf("expected runtime error, got '%v'", err)
	}
}

func expectResult(t *testing.T, src string, expected interface{}) {
	t.Helper()
	if result, err := interpret(t, src); err != nil {
		t.Error(err)
	} else {
		if result != expected {
			t.Errorf("expected '%v', got '%v'", expected, result)
		}
	}
}

func interpret(t *testing.T, src string) (interface{}, error) {
	expr, _ := parser.NewParser(scanner.NewScanner(bufio.NewReader(strings.NewReader(src)))).Parse()
	return NewInterpreter().Interpret(expr)
}
