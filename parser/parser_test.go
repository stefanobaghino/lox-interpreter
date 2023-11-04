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

func TestParserWrongTokenAfterParen(t *testing.T) {
	expectErrors(t, "(1 + 2var", "expected '\\)' after expression \\(at 'var'\\)")
}

func TestParserUnclosedParen(t *testing.T) {
	expectErrors(t, "(1 + 2", "expected '\\)' after expression \\(at end\\)")
}

func TestParserLexicalErrors(t *testing.T) {
	expectErrors(t, "(1 + 2%", "unexpected character")
}

func TestParserMissingExpression(t *testing.T) {
	expectErrors(t, "(", "expected expression")
}

func TestParserMissingExpressionFromStatement(t *testing.T) {
	expectErrors(t, "print;", "expected expression")
}

func TestParserVarDecl(t *testing.T) {
	expectFormatted(t, "var x = 1;")
}

func TestParserInvalidAssignmentTarget(t *testing.T) {
	expectErrors(t, "1 = 2;", "invalid assignment target")
}

func TestParserMultiErr(t *testing.T) {
	expectErrors(t, "1 = 2; print;", "invalid assignment target", "expected expression")
}

func TestParserAssert(t *testing.T) {
	expectFormatted(t, "assert true;")
}

func TestParserStatements(t *testing.T) {
	expectFormatted(t, "print 1;\n{\n\tvar x = 1;\n\tx = 2;\n}")
}

func TestParserIf(t *testing.T) {
	expectFormatted(t, "if (true)\n{\n\tprint 1;\n}")
}

func TestParserCall(t *testing.T) {
	expectFormatted(t, "foo(1);")
	expectFormatted(t, "curried(1)(2);")
}

func TestParserFunDecl(t *testing.T) {
	expectFormatted(t, "fun foo()\n{\n\tprint \"a\";\n\tprint clock();\n}")
}

func expectErrors(t *testing.T, src string, regexps ...string) {
	t.Helper()
	p := NewParser(scanner.NewScanner(bufio.NewReader(strings.NewReader(src))))
	for {
		if stmt, err := p.NextStatement(); err != nil {
			if len(regexps) == 0 {
				t.Fatalf("unexpected error: %v", err)
			}
			re := regexps[0]
			regexps = regexps[1:]
			if !regexp.MustCompile(re).MatchString(err.Error()) {
				t.Errorf("expected '%s' to match '%v'", re, err.Error())
			}
		} else if _, ok := stmt.(*ast.EndStmt); ok {
			break
		}
	}
	if len(regexps) > 0 {
		t.Errorf("expected '%v' errors, got none", regexps)
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
