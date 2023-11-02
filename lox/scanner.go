package lox

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type Scanner struct {
	reader  *bufio.Reader
	chars   []rune
	current int
	line    int
}

type syntaxError struct {
	line    int
	message string
}

func (e syntaxError) Error() string {
	return fmt.Sprintf("syntax error on line %d: %s", e.line, e.message)
}

func NewScanner(reader *bufio.Reader) *Scanner {
	return &Scanner{reader: reader, line: 1}
}

func (s *Scanner) NextToken() (Token, error) {
next:
	s.chars = s.chars[s.current:]
	s.current = 0
	r := s.advance()
	switch {
	case r == utf8.RuneError:
		return s.mkToken(EOF), nil
	case r == '(':
		return s.mkToken(LEFT_PAREN), nil
	case r == ')':
		return s.mkToken(RIGHT_PAREN), nil
	case r == '{':
		return s.mkToken(LEFT_BRACE), nil
	case r == '}':
		return s.mkToken(RIGHT_BRACE), nil
	case r == ',':
		return s.mkToken(COMMA), nil
	case r == '.':
		return s.mkToken(DOT), nil
	case r == '-':
		return s.mkToken(MINUS), nil
	case r == '+':
		return s.mkToken(PLUS), nil
	case r == ';':
		return s.mkToken(SEMICOLON), nil
	case r == '*':
		return s.mkToken(STAR), nil
	case r == '!':
		if s.match('=') {
			return s.mkToken(BANG_EQUAL), nil
		} else {
			return s.mkToken(BANG), nil
		}
	case r == '=':
		if s.match('=') {
			return s.mkToken(EQUAL_EQUAL), nil
		} else {
			return s.mkToken(EQUAL), nil
		}
	case r == '<':
		if s.match('=') {
			return s.mkToken(LESS_EQUAL), nil
		} else {
			return s.mkToken(LESS), nil
		}
	case r == '>':
		if s.match('=') {
			return s.mkToken(GREATER_EQUAL), nil
		} else {
			return s.mkToken(GREATER), nil
		}
	case r == '/':
		if s.match('/') {
			s.skipUntil(func(r rune) bool { return r == '\n' })
			goto next
		} else {
			return s.mkToken(SLASH), nil
		}
	case unicode.IsSpace(r):
		if r == '\n' {
			s.line += 1
		}
		goto next
	case r == '"':
		return s.str(), nil
	case unicode.IsDigit(r):
		return s.num()
	case unicode.IsLetter(r):
		return s.id(), nil
	default:
		return s.mkToken(ERROR), &syntaxError{s.line, "Unexpected character."}
	}
}

func (s *Scanner) id() Token {
	s.skipUntil(func(r rune) bool { return !unicode.IsDigit(r) && !unicode.IsLetter(r) })
	t, ok := keywords[string(s.chars[:s.current])]
	if !ok {
		t = IDENTIFIER
	}
	return s.mkToken(t)
}

func (s *Scanner) num() (Token, error) {
	s.skipUntil(func(r rune) bool { return !unicode.IsDigit(r) })
	// Look for a fractional part.
	if unicode.IsDigit(s.readRuneAhead(1)) && s.readRune() == '.' {
		s.current += 1 // skip the dot
		s.skipUntil(func(r rune) bool { return !unicode.IsDigit(r) })
	}
	if x, err := strconv.ParseFloat(string(s.chars[:s.current]), 64); err != nil {
		return s.mkToken(ERROR), &syntaxError{s.line, "Invalid number."}
	} else {
		return s.mkLiteral(NUMBER, x), nil
	}
}

func (s *Scanner) str() Token {
	s.current += 1 // skip the opening quote
	s.skipUntil(func(r rune) bool { return r == '"' })
	value := string(s.chars[1:s.current])
	s.current += 1 // skip the closing quote
	return s.mkLiteral(STRING, value)
}

func (s *Scanner) match(expected rune) bool {
	r := s.readRune()
	if r != expected {
		return false
	}
	s.current += 1
	return true
}

func (s *Scanner) advance() rune {
	r := s.readRune()
	if r != utf8.RuneError {
		s.current += 1
	}
	return r
}

func (s *Scanner) skipUntil(p func(rune) bool) {
	for {
		r := s.readRune()
		if r == utf8.RuneError || p(r) {
			break
		}
		s.current += 1
	}
}

func (s *Scanner) readRune() rune {
	return s.readRuneAhead(0)
}

func (s *Scanner) readRuneAhead(offset int) rune {
	offset += s.current
	for d := offset - len(s.chars) + 1; d > 0; d-- {
		r, sz, err := s.reader.ReadRune()
		if r == utf8.RuneError && sz == 1 {
			panic(errors.New("invalid UTF-8 sequence"))
		}
		if err != nil {
			if err == io.EOF {
				s.chars = append(s.chars, utf8.RuneError)
			} else {
				panic(err)
			}
		}
		s.chars = append(s.chars, r)
	}
	return s.chars[offset]
}

func (s *Scanner) mkToken(t TokenType) Token {
	return s.mkLiteral(t, nil)
}

func (s *Scanner) mkLiteral(t TokenType, literal interface{}) Token {
	lexeme := string(s.chars[:s.current])
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
