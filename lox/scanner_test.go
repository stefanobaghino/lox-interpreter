package lox

import (
	"bufio"
	"strings"
	"testing"
)

func TestScannerEmpty(t *testing.T) {
	src := ""
	s := NewScanner(bufio.NewReader(strings.NewReader(src)))
	expectTokenType(t, expectNext(t, s), EOF)
}

func TestScannerError(t *testing.T) {
	src := "#"
	s := NewScanner(bufio.NewReader(strings.NewReader(src)))
	expectErrorMessage(t, expectSyntaxError(t, s), "Unexpected character.")
	expectTokenType(t, expectNext(t, s), EOF)
}

func TestScannerErrorKeepsGoing(t *testing.T) {
	src := "#\n#kbye"
	s := NewScanner(bufio.NewReader(strings.NewReader(src)))
	expectErrorMessage(t, expectSyntaxError(t, s), "Unexpected character.")
	expectErrorMessage(t, expectSyntaxError(t, s), "Unexpected character.")
	expectIdentifier(t, expectNext(t, s), "kbye")
	expectTokenType(t, expectNext(t, s), EOF)
}

func TestScannerNumberTooBig(t *testing.T) {
	src := strings.Builder{}
	for i := 0; i < 1000; i++ {
		src.WriteRune('9')
	}
	s := NewScanner(bufio.NewReader(strings.NewReader(src.String())))
	expectErrorMessage(t, expectSyntaxError(t, s), "Invalid number.")
}

func TestScannerSimpleTokens(t *testing.T) {
	src := "(){},.+-;*"
	s := NewScanner(bufio.NewReader(strings.NewReader(src)))
	expectTokenType(t, expectNext(t, s), LEFT_PAREN)
	expectTokenType(t, expectNext(t, s), RIGHT_PAREN)
	expectTokenType(t, expectNext(t, s), LEFT_BRACE)
	expectTokenType(t, expectNext(t, s), RIGHT_BRACE)
	expectTokenType(t, expectNext(t, s), COMMA)
	expectTokenType(t, expectNext(t, s), DOT)
	expectTokenType(t, expectNext(t, s), PLUS)
	expectTokenType(t, expectNext(t, s), MINUS)
	expectTokenType(t, expectNext(t, s), SEMICOLON)
	expectTokenType(t, expectNext(t, s), STAR)
	expectTokenType(t, expectNext(t, s), EOF)
}

func TestScannerWhiteSpaceAndComplexTokens(t *testing.T) {
	src := `
      == != <= >= // comment
      = ! < > // another comment
      1	/	2 // and another one`
	s := NewScanner(bufio.NewReader(strings.NewReader(src)))
	expectTokenType(t, expectNext(t, s), EQUAL_EQUAL)
	expectTokenType(t, expectNext(t, s), BANG_EQUAL)
	expectTokenType(t, expectNext(t, s), LESS_EQUAL)
	expectTokenType(t, expectNext(t, s), GREATER_EQUAL)
	expectTokenType(t, expectNext(t, s), EQUAL)
	expectTokenType(t, expectNext(t, s), BANG)
	expectTokenType(t, expectNext(t, s), LESS)
	expectTokenType(t, expectNext(t, s), GREATER)
	expectTokenType(t, expectNext(t, s), NUMBER)
	expectTokenType(t, expectNext(t, s), SLASH)
	expectTokenType(t, expectNext(t, s), NUMBER)
	expectTokenType(t, expectNext(t, s), EOF)
}

func TestScannerNumbers(t *testing.T) {
	src := "123 123.456 0.456"
	s := NewScanner(bufio.NewReader(strings.NewReader(src)))
	expectNumberLiteral(t, expectNext(t, s), 123)
	expectNumberLiteral(t, expectNext(t, s), 123.456)
	expectNumberLiteral(t, expectNext(t, s), 0.456)
	expectTokenType(t, expectNext(t, s), EOF)
}

func TestScannerStrings(t *testing.T) {
	src := `"hello" "world"`
	s := NewScanner(bufio.NewReader(strings.NewReader(src)))
	expectStringLiteral(t, expectNext(t, s), "hello")
	expectStringLiteral(t, expectNext(t, s), "world")
	expectTokenType(t, expectNext(t, s), EOF)
}

func TestScannerIdentifiers(t *testing.T) {
	src := `hello world`
	s := NewScanner(bufio.NewReader(strings.NewReader(src)))
	expectIdentifier(t, expectNext(t, s), "hello")
	expectIdentifier(t, expectNext(t, s), "world")
	expectTokenType(t, expectNext(t, s), EOF)
}

func expectSyntaxError(t *testing.T, scanner *Scanner) *syntaxError {
	t.Helper()
	r, err := scanner.NextToken()
	if err == nil {
		t.Fatalf("expected error, got '%v'", r)
	}
	syntaxErr, ok := err.(*syntaxError)
	if !ok {
		t.Fatalf("expected syntaxError, got '%T'", err)
	}
	return syntaxErr
}

func expectErrorMessage(t *testing.T, err *syntaxError, expected string) {
	t.Helper()
	if err.message != expected {
		t.Errorf("expected message '%s', got '%s'", expected, err.message)
	}
}

func expectNext(t *testing.T, scanner *Scanner) Token {
	t.Helper()
	token, err := scanner.NextToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return token
}

func expectTokenType(t *testing.T, token Token, expected TokenType) {
	t.Helper()
	if token.Type != expected {
		t.Errorf("expected token type %s, got %s", expected, token.Type)
	}
}

func expectNumberLiteral(t *testing.T, token Token, expected float64) {
	t.Helper()
	expectTokenType(t, token, NUMBER)
	if token.Literal != expected {
		t.Errorf("expected number literal %f, got %v", expected, token.Literal)
	}
}

func expectStringLiteral(t *testing.T, token Token, expected string) {
	t.Helper()
	expectTokenType(t, token, STRING)
	if token.Literal != expected {
		t.Errorf("expected string literal %s, got %v", expected, token.Literal)
	}
}

func expectIdentifier(t *testing.T, token Token, expected string) {
	t.Helper()
	expectTokenType(t, token, IDENTIFIER)
	if token.Lexeme != expected {
		t.Errorf("expected identifier %s, got %v", expected, token.Lexeme)
	}
}
