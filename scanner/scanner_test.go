package scanner

import (
	"bufio"
	"lox/token"
	"strings"
	"testing"
)

func TestScannerEmpty(t *testing.T) {
	src := ""
	s := NewScanner(bufio.NewReader(strings.NewReader(src)))
	expectTokenType(t, expectNext(t, s), token.EOF)
}

func TestScannerError(t *testing.T) {
	src := "#"
	s := NewScanner(bufio.NewReader(strings.NewReader(src)))
	expectErrorMessage(t, expectLexicalError(t, s), "unexpected character")
	expectTokenType(t, expectNext(t, s), token.EOF)
}

func TestScannerErrorKeepsGoing(t *testing.T) {
	src := "#\n#kbye"
	s := NewScanner(bufio.NewReader(strings.NewReader(src)))
	expectErrorMessage(t, expectLexicalError(t, s), "unexpected character")
	expectErrorMessage(t, expectLexicalError(t, s), "unexpected character")
	expectIdentifier(t, expectNext(t, s), "kbye")
	expectTokenType(t, expectNext(t, s), token.EOF)
}

func TestScannerNumberTooBig(t *testing.T) {
	src := strings.Builder{}
	for i := 0; i < 1000; i++ {
		src.WriteRune('9')
	}
	s := NewScanner(bufio.NewReader(strings.NewReader(src.String())))
	expectErrorMessage(t, expectLexicalError(t, s), "invalid number")
}

func TestScannerSimpleTokens(t *testing.T) {
	src := "(){},.+-;*"
	s := NewScanner(bufio.NewReader(strings.NewReader(src)))
	expectTokenType(t, expectNext(t, s), token.LEFT_PAREN)
	expectTokenType(t, expectNext(t, s), token.RIGHT_PAREN)
	expectTokenType(t, expectNext(t, s), token.LEFT_BRACE)
	expectTokenType(t, expectNext(t, s), token.RIGHT_BRACE)
	expectTokenType(t, expectNext(t, s), token.COMMA)
	expectTokenType(t, expectNext(t, s), token.DOT)
	expectTokenType(t, expectNext(t, s), token.PLUS)
	expectTokenType(t, expectNext(t, s), token.MINUS)
	expectTokenType(t, expectNext(t, s), token.SEMICOLON)
	expectTokenType(t, expectNext(t, s), token.STAR)
	expectTokenType(t, expectNext(t, s), token.EOF)
}

func TestScannerWhiteSpaceAndComplexTokens(t *testing.T) {
	src := `
	== != <= >=	// comment
	= ! < >		// another comment
	1	/	2	// and another one`
	s := NewScanner(bufio.NewReader(strings.NewReader(src)))
	expectTokenType(t, expectNext(t, s), token.EQUAL_EQUAL)
	expectTokenType(t, expectNext(t, s), token.BANG_EQUAL)
	expectTokenType(t, expectNext(t, s), token.LESS_EQUAL)
	expectTokenType(t, expectNext(t, s), token.GREATER_EQUAL)
	expectTokenType(t, expectNext(t, s), token.EQUAL)
	expectTokenType(t, expectNext(t, s), token.BANG)
	expectTokenType(t, expectNext(t, s), token.LESS)
	expectTokenType(t, expectNext(t, s), token.GREATER)
	expectTokenType(t, expectNext(t, s), token.NUMBER)
	expectTokenType(t, expectNext(t, s), token.SLASH)
	expectTokenType(t, expectNext(t, s), token.NUMBER)
	expectTokenType(t, expectNext(t, s), token.EOF)
}

func TestScannerNumbers(t *testing.T) {
	src := "123 123.456 0.456"
	s := NewScanner(bufio.NewReader(strings.NewReader(src)))
	expectNumberLiteral(t, expectNext(t, s), 123)
	expectNumberLiteral(t, expectNext(t, s), 123.456)
	expectNumberLiteral(t, expectNext(t, s), 0.456)
	expectTokenType(t, expectNext(t, s), token.EOF)
}

func TestScannerStrings(t *testing.T) {
	src := `"hello" "world"`
	s := NewScanner(bufio.NewReader(strings.NewReader(src)))
	expectStringLiteral(t, expectNext(t, s), "hello")
	expectStringLiteral(t, expectNext(t, s), "world")
	expectTokenType(t, expectNext(t, s), token.EOF)
}

func TestScannerIdentifiers(t *testing.T) {
	src := `hello world`
	s := NewScanner(bufio.NewReader(strings.NewReader(src)))
	expectIdentifier(t, expectNext(t, s), "hello")
	expectIdentifier(t, expectNext(t, s), "world")
	expectTokenType(t, expectNext(t, s), token.EOF)
}

func TestScannerLineNumbers(t *testing.T) {
	src := `
	2 2	

	4	4 
	5

	7
	`
	s := NewScanner(bufio.NewReader(strings.NewReader(src)))
	expectLineNumber(t, expectNext(t, s), 2)
	expectLineNumber(t, expectNext(t, s), 2)
	expectLineNumber(t, expectNext(t, s), 4)
	expectLineNumber(t, expectNext(t, s), 4)
	expectLineNumber(t, expectNext(t, s), 5)
	expectLineNumber(t, expectNext(t, s), 7)
	expectTokenType(t, expectNext(t, s), token.EOF)
}

func expectLexicalError(t *testing.T, scanner *Scanner) *LexicalError {
	t.Helper()
	r, err := scanner.NextToken()
	if err == nil {
		t.Fatalf("expected error, got '%v'", r)
	}
	lexicalErr, ok := err.(*LexicalError)
	if !ok {
		t.Fatalf("expected syntaxError, got '%T'", err)
	}
	return lexicalErr
}

func expectErrorMessage(t *testing.T, err *LexicalError, expected string) {
	t.Helper()
	if err.message != expected {
		t.Errorf("expected message '%s', got '%s'", expected, err.message)
	}
}

func expectNext(t *testing.T, scanner *Scanner) token.Token {
	t.Helper()
	token, err := scanner.NextToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return token
}

func expectTokenType(t *testing.T, token token.Token, expected token.Type) {
	t.Helper()
	if token.Type != expected {
		t.Errorf("expected token type %s, got %s", expected, token.Type)
	}
}

func expectNumberLiteral(t *testing.T, tk token.Token, expected float64) {
	t.Helper()
	expectTokenType(t, tk, token.NUMBER)
	if tk.Literal != expected {
		t.Errorf("expected number literal %f, got %v", expected, tk.Literal)
	}
}

func expectStringLiteral(t *testing.T, tk token.Token, expected string) {
	t.Helper()
	expectTokenType(t, tk, token.STRING)
	if tk.Literal != expected {
		t.Errorf("expected string literal %s, got %v", expected, tk.Literal)
	}
}

func expectIdentifier(t *testing.T, tk token.Token, expected string) {
	t.Helper()
	expectTokenType(t, tk, token.IDENTIFIER)
	if tk.Lexeme != expected {
		t.Errorf("expected identifier %s, got %v", expected, tk.Lexeme)
	}
}

func expectLineNumber(t *testing.T, tk token.Token, expected int) {
	t.Helper()
	if tk.Line != expected {
		t.Errorf("expected line number %d, got %d", expected, tk.Line)
	}
}
