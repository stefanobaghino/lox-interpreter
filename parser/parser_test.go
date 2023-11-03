package parser

import (
	"bufio"
	"lox/ast"
	"lox/format"
	"lox/scanner"
	"regexp"
	"strings"
	"testing"
)

func TestParserOne(t *testing.T) {
	expectFormatted(t, "1;")
}

func TestParserBinary(t *testing.T) {
	expectFormatted(t, "1 + 2;")
}

func TestParserUnary(t *testing.T) {
	expectFormatted(t, "-1;")
}

func TestParserUnaryAndBinary(t *testing.T) {
	expectFormatted(t, "-1 * -1;")
}

func TestParserGrouping(t *testing.T) {
	expectFormatted(t, "1 * (2 + 3);")
}

func TestParserBinaryAssociativity(t *testing.T) {
	expectFormatted(t, "1 + 2 + 3;")
}

func TestParserBooleans(t *testing.T) {
	expectFormatted(t, "true == !false;")
}

func TestParserComparisons(t *testing.T) {
	expectFormatted(t, "0 <= 1 >= 2;")
}

func TestParserNil(t *testing.T) {
	expectFormatted(t, "nil - nil == 0;")
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

func TestParserMissingExpression(t *testing.T) {
	expectError(t, "(", "expected expression")
}

func TestParserVarDecl(t *testing.T) {
	expectFormatted(t, "var x = 1;")
}

func TestParserAssert(t *testing.T) {
	expectFormatted(t, "assert true;")
}

func TestParserStatements(t *testing.T) {
	expectFormatted(t, "print 1;\n{\n\tvar x = 1;\n\tx = 2;\n}")
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

func expectFormatted(t *testing.T, src string) {
	t.Helper()
	p := NewParser(scanner.NewScanner(bufio.NewReader(strings.NewReader(src))))
	f := format.NewFormatter()
	builder := strings.Builder{}
	for {
		if stmt, err := p.NextStatement(); err != nil {
			t.Error(err)
		} else if _, ok := stmt.(*ast.EndStmt); !ok {
			builder.WriteString(f.Format(stmt))
		} else {
			break
		}
	}
	result := builder.String()
	if result != src {
		t.Errorf("expected '%s', got '%s'", src, result)
	}
}
