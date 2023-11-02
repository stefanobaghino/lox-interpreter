package lox

import (
	"bufio"
	"regexp"
	"strings"
	"testing"
)

func TestParserOne(t *testing.T) {
	src := "1"
	expectExpression(t, parser(src), "1")
}

func TestParserBinary(t *testing.T) {
	src := "1 + 2"
	expectExpression(t, parser(src), "(+ 1 2)")
}

func TestParserUnary(t *testing.T) {
	src := "-1"
	expectExpression(t, parser(src), "(- 1)")
}

func TestParserUnaryAndBinary(t *testing.T) {
	src := "-1 * -1"
	expectExpression(t, parser(src), "(* (- 1) (- 1))")
}

func TestParserGrouping(t *testing.T) {
	src := "1 * (2 + 3)"
	expectExpression(t, parser(src), "(* 1 (group (+ 2 3)))")
}

func TestParserBinaryAssociativity(t *testing.T) {
	src := "1 + 2 + 3"
	expectExpression(t, parser(src), "(+ (+ 1 2) 3)")
}

func TestParserBooleans(t *testing.T) {
	src := "true == !false"
	expectExpression(t, parser(src), "(== true (! false))")
}

func TestParserComparisons(t *testing.T) {
	src := "0 <= 1 >= 2"
	expectExpression(t, parser(src), "(>= (<= 0 1) 2)")
}

func TestParserNil(t *testing.T) {
	src := "nil - nil == 0"
	expectExpression(t, parser(src), "(== (- nil nil) 0)")
}

func TestParserEmpty(t *testing.T) {
	src := ""
	expectError(t, parser(src), "expected expression")
}

func TestParserWrongParen(t *testing.T) {
	src := "(1 + 2("
	expectError(t, parser(src), "expected '\\)' after expression \\(at '\\('\\)")
}

func TestParserUnclosedParen(t *testing.T) {
	src := "(1 + 2"
	expectError(t, parser(src), "expected '\\)' after expression \\(at end\\)")
}

func TestParserLexicalErrors(t *testing.T) {
	src := "(1 + 2%"
	expectError(t, parser(src), "unexpected character")
}

func expectError(t *testing.T, p *Parser, re string) {
	t.Helper()
	if _, err := p.Parse(); err == nil {
		t.Errorf("expected error, got none")
	} else {
		if !regexp.MustCompile(re).MatchString(err.Error()) {
			t.Errorf("expected '%s' to match '%v'", re, err.Error())
		}
	}
}

func expectExpression(t *testing.T, p *Parser, expected string) {
	t.Helper()
	if expr, err := p.Parse(); err != nil {
		t.Error(err)
	} else {
		result := expr.Accept(AstPrinter).(string)
		if result != expected {
			t.Errorf("expected '%s', got '%s'", expected, result)
		}
	}
}

func parser(src string) *Parser {
	return NewParser(NewScanner(bufio.NewReader(strings.NewReader(src))))
}
