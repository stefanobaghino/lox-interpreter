package parser

import (
	"bufio"
	"lox/ast"
	"lox/scanner"
	"regexp"
	"strings"
	"testing"
)

func TestParserOne(t *testing.T) {
	expectExpression(t, "1;", "1")
}

func TestParserBinary(t *testing.T) {
	expectExpression(t, "1 + 2;", "(+ 1 2)")
}

func TestParserUnary(t *testing.T) {
	expectExpression(t, "-1;", "(- 1)")
}

func TestParserUnaryAndBinary(t *testing.T) {
	expectExpression(t, "-1 * -1;", "(* (- 1) (- 1))")
}

func TestParserGrouping(t *testing.T) {
	expectExpression(t, "1 * (2 + 3);", "(* 1 (group (+ 2 3)))")
}

func TestParserBinaryAssociativity(t *testing.T) {
	expectExpression(t, "1 + 2 + 3;", "(+ (+ 1 2) 3)")
}

func TestParserBooleans(t *testing.T) {
	expectExpression(t, "true == !false;", "(== true (! false))")
}

func TestParserComparisons(t *testing.T) {
	expectExpression(t, "0 <= 1 >= 2;", "(>= (<= 0 1) 2)")
}

func TestParserNil(t *testing.T) {
	expectExpression(t, "nil - nil == 0;", "(== (- nil nil) 0)")
}

func TestParserWrongParen(t *testing.T) {
	expectError(t, "(1 + 2(", "expected '\\)' after expression \\(at '\\('\\)")
}

func TestParserUnclosedParen(t *testing.T) {
	expectError(t, "(1 + 2", "expected '\\)' after expression \\(at end\\)")
}

func TestParserLexicalErrors(t *testing.T) {
	expectError(t, "(1 + 2%", "unexpected character")
}

func expectError(t *testing.T, src string, re string) {
	t.Helper()
	p := NewParser(scanner.NewScanner(bufio.NewReader(strings.NewReader(src))))
	if _, err := p.NextStatement(); err == nil {
		t.Errorf("expected error, got none")
	} else {
		if !regexp.MustCompile(re).MatchString(err.Error()) {
			t.Errorf("expected '%s' to match '%v'", re, err.Error())
		}
	}
}

func expectExpression(t *testing.T, src string, expected string) {
	t.Helper()
	p := NewParser(scanner.NewScanner(bufio.NewReader(strings.NewReader(src))))
	if stmt, err := p.NextStatement(); err != nil {
		t.Error(err)
	} else if expr, ok := stmt.(*ast.ExprStmt); !ok {
		t.Errorf("expected expression, got '%v'", stmt)
	} else {
		result := ast.Print(expr.Expression)
		if result != expected {
			t.Errorf("expected '%s', got '%s'", expected, result)
		}
	}
}
