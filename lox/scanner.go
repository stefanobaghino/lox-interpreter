package lox

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type Scanner struct {
	source  []byte
	start   int
	current int
	line    int
	broken  bool
}

type syntaxError struct {
	line    int
	message string
}

func (e syntaxError) Error() string {
	return fmt.Sprintf("syntax error on line %d: %s", e.line, e.message)
}

func NewScanner(source []byte) *Scanner {
	return &Scanner{source: source, line: 1}
}

func (s *Scanner) NextToken() (Token, error) {
next:
	s.start = s.current

	r, err := s.advance()
	if r == utf8.RuneError {
		return s.mkToken(EOF), err
	}
	switch r {
	case '(':
		return s.mkToken(LEFT_PAREN), nil
	case ')':
		return s.mkToken(RIGHT_PAREN), nil
	case '{':
		return s.mkToken(LEFT_BRACE), nil
	case '}':
		return s.mkToken(RIGHT_BRACE), nil
	case ',':
		return s.mkToken(COMMA), nil
	case '.':
		return s.mkToken(DOT), nil
	case '-':
		return s.mkToken(MINUS), nil
	case '+':
		return s.mkToken(PLUS), nil
	case ';':
		return s.mkToken(SEMICOLON), nil
	case '*':
		return s.mkToken(STAR), nil
	case '!':
		match, err := s.match('=')
		if err != nil {
			return s.mkToken(ERROR), err
		}
		if match {
			return s.mkToken(BANG_EQUAL), nil
		} else {
			return s.mkToken(BANG), nil
		}
	case '=':
		match, err := s.match('=')
		if err != nil {
			return s.mkToken(ERROR), err
		}
		if match {
			return s.mkToken(EQUAL_EQUAL), nil
		} else {
			return s.mkToken(EQUAL), nil
		}
	case '<':
		match, err := s.match('=')
		if err != nil {
			return s.mkToken(ERROR), err
		}
		if match {
			return s.mkToken(LESS_EQUAL), nil
		} else {
			return s.mkToken(LESS), nil
		}
	case '>':
		match, err := s.match('=')
		if err != nil {
			return s.mkToken(ERROR), err
		}
		if match {
			return s.mkToken(GREATER_EQUAL), nil
		} else {
			return s.mkToken(GREATER), nil
		}
	case '/':
		match, err := s.match('/')
		if err != nil {
			return s.mkToken(ERROR), err
		}
		if match {
			s.skipUntil(func(r rune) bool { return r == '\n' })
			goto next
		} else {
			return s.mkToken(SLASH), nil
		}
	case '\n':
		s.line++
		goto next
	case ' ':
		goto next
	case '\r':
		goto next
	case '\t':
		goto next
	case '"':
		return s.str()
	default:
		if unicode.IsDigit(r) {
			return s.num()
		} else if unicode.IsLetter(r) {
			return s.id()
		}
	}
	return s.mkToken(ERROR), &syntaxError{s.line, "Unexpected character."}
}

func (s *Scanner) id() (Token, error) {
	if err := s.skipUntil(func(r rune) bool { return !unicode.IsDigit(r) && !unicode.IsLetter(r) }); err != nil {
		return s.mkToken(ERROR), err
	}
	t, ok := keywords[string(s.source[s.start:s.current])]
	if !ok {
		t = IDENTIFIER
	}
	return s.mkToken(t), nil
}

func (s *Scanner) num() (Token, error) {
	if err := s.skipUntil(func(r rune) bool { return !unicode.IsDigit(r) }); err != nil {
		return s.mkToken(ERROR), err
	}
	// Look for a fractional part.
	if rs, _, err := s.lookahead(2); err != nil {
		return s.mkToken(ERROR), err
	} else if len(rs) == 2 && rs[0] == '.' && unicode.IsDigit(rs[1]) {
		s.current += 1 // '.' is one byte, move past it
		if err := s.skipUntil(func(r rune) bool { return !unicode.IsDigit(r) }); err != nil {
			return s.mkToken(ERROR), err
		}
	}
	if x, err := strconv.ParseFloat(string(s.source[s.start:s.current]), 64); err != nil {

		return s.mkToken(ERROR), &syntaxError{s.line, "Invalid number."}
	} else {
		return s.mkLiteral(NUMBER, x), nil
	}
}

func (s *Scanner) str() (Token, error) {
	s.current += 1 // '"' is one byte, move past it
	if err := s.skipUntil(func(r rune) bool { return r == '"' }); err != nil {
		return s.mkToken(ERROR), err
	}
	value := string(s.source[s.start+1 : s.current])
	s.current += 1 // skip the closing quote
	return s.mkLiteral(STRING, value), nil
}

func (s *Scanner) match(expected rune) (bool, error) {
	r, sz, err := s.peek()
	if err != nil {
		return false, err
	}
	if r != expected {
		return false, nil
	}
	s.current += sz
	return true, nil
}

func (s *Scanner) advance() (rune, error) {
	r, sz, err := s.peek()
	if err != nil {
		return r, err
	}
	s.current += sz
	return r, nil
}

func (s *Scanner) skipUntil(p func(rune) bool) error {
	for {
		r, sz, err := s.peek()
		if err != nil {
			return err
		}
		if sz == 0 {
			break
		}
		if p(r) {
			break
		}
		s.current += sz
	}
	return nil
}

func (s *Scanner) peek() (rune, int, error) {
	return s.peekAhead(0)
}

func (s *Scanner) lookahead(n int) ([]rune, int, error) {
	rs := make([]rune, 0, n)
	read := 0
	for i := 0; i < n; i++ {
		r, sz, err := s.peekAhead(read)
		read += sz
		if err != nil || sz == 0 {
			return nil, read, err
		}
		rs = append(rs, r)
	}
	return rs, read, nil
}

func (s *Scanner) peekAhead(offset int) (rune, int, error) {
	r, sz := utf8.DecodeRune(s.source[s.current+offset:])
	if r == utf8.RuneError && sz == 1 {
		s.broken = true
		return r, 1, &syntaxError{s.line, "Invalid UTF-8 sequence."}
	}
	return r, sz, nil
}

func (s *Scanner) mkToken(t TokenType) Token {
	return s.mkLiteral(t, nil)
}

func (s *Scanner) mkLiteral(t TokenType, literal interface{}) Token {
	lexeme := string(s.source[s.start:s.current])
	return Token{Type: t, Lexeme: lexeme, Literal: literal, Line: s.line}
}

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}
