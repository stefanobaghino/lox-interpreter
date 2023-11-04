package scanner

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"lox/token"
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

type LexicalError struct {
	line    int
	message string
}

func (e LexicalError) Error() string {
	return fmt.Sprintf("lexical error on line %d: %s", e.line, e.message)
}

func (e LexicalError) Line() int {
	return e.line
}

func NewScanner(reader *bufio.Reader) *Scanner {
	return &Scanner{reader: reader, line: 1}
}

func (s *Scanner) NextToken() (token.Token, error) {
	s.chars = s.chars[s.current:]
	s.current = 0
	r := s.advance()
	switch {
	case r == utf8.RuneError:
		return s.mkToken(token.EOF), nil
	case r == '(':
		return s.mkToken(token.LEFT_PAREN), nil
	case r == ')':
		return s.mkToken(token.RIGHT_PAREN), nil
	case r == '{':
		return s.mkToken(token.LEFT_BRACE), nil
	case r == '}':
		return s.mkToken(token.RIGHT_BRACE), nil
	case r == ',':
		return s.mkToken(token.COMMA), nil
	case r == '.':
		return s.mkToken(token.DOT), nil
	case r == '-':
		return s.mkToken(token.MINUS), nil
	case r == '+':
		return s.mkToken(token.PLUS), nil
	case r == ';':
		return s.mkToken(token.SEMICOLON), nil
	case r == '*':
		return s.mkToken(token.STAR), nil
	case r == '!':
		if s.match('=') {
			return s.mkToken(token.BANG_EQUAL), nil
		} else {
			return s.mkToken(token.BANG), nil
		}
	case r == '=':
		if s.match('=') {
			return s.mkToken(token.EQUAL_EQUAL), nil
		} else {
			return s.mkToken(token.EQUAL), nil
		}
	case r == '<':
		if s.match('=') {
			return s.mkToken(token.LESS_EQUAL), nil
		} else {
			return s.mkToken(token.LESS), nil
		}
	case r == '>':
		if s.match('=') {
			return s.mkToken(token.GREATER_EQUAL), nil
		} else {
			return s.mkToken(token.GREATER), nil
		}
	case r == '/':
		if s.match('/') {
			s.skipUntil(func(r rune) bool { return r == '\n' })
			return s.NextToken()
		} else {
			return s.mkToken(token.SLASH), nil
		}
	case unicode.IsSpace(r):
		if r == '\n' {
			s.line += 1
		}
		s.skipUntil(func(r rune) bool {
			if r == '\n' {
				s.line += 1
			}
			return !unicode.IsSpace(r)
		})
		return s.NextToken()
	case r == '"':
		return s.str(), nil
	case unicode.IsDigit(r):
		return s.num()
	case unicode.IsLetter(r):
		return s.id(), nil
	default:
		return s.mkToken(token.ERROR), &LexicalError{s.line, "unexpected character"}
	}
}

func (s *Scanner) id() token.Token {
	s.skipUntil(func(r rune) bool { return !unicode.IsDigit(r) && !unicode.IsLetter(r) })
	t, ok := keywords[string(s.chars[:s.current])]
	if !ok {
		t = token.IDENTIFIER
	}
	return s.mkToken(t)
}

func (s *Scanner) num() (token.Token, error) {
	s.skipUntil(func(r rune) bool { return !unicode.IsDigit(r) })
	// Look for a fractional part.
	if unicode.IsDigit(s.readRuneAhead(1)) && s.readRune() == '.' {
		s.current += 1 // skip the dot
		s.skipUntil(func(r rune) bool { return !unicode.IsDigit(r) })
	}
	if x, err := strconv.ParseFloat(string(s.chars[:s.current]), 64); err != nil {
		return s.mkToken(token.ERROR), &LexicalError{s.line, "invalid number"}
	} else {
		return s.mkLiteral(token.NUMBER, x), nil
	}
}

func (s *Scanner) str() token.Token {
	s.current += 1 // skip the opening quote
	s.skipUntil(func(r rune) bool { return r == '"' })
	value := string(s.chars[1:s.current])
	s.current += 1 // skip the closing quote
	return s.mkLiteral(token.STRING, value)
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

func (s *Scanner) mkToken(t token.Type) token.Token {
	return s.mkLiteral(t, nil)
}

func (s *Scanner) mkLiteral(t token.Type, literal interface{}) token.Token {
	lexeme := string(s.chars[:s.current])
	return token.Token{Type: t, Lexeme: lexeme, Literal: literal, Line: s.line}
}

var keywords = map[string]token.Type{
	"and":    token.AND,
	"assert": token.ASSERT,
	"class":  token.CLASS,
	"else":   token.ELSE,
	"false":  token.FALSE,
	"for":    token.FOR,
	"fun":    token.FUN,
	"if":     token.IF,
	"nil":    token.NIL,
	"or":     token.OR,
	"print":  token.PRINT,
	"return": token.RETURN,
	"super":  token.SUPER,
	"this":   token.THIS,
	"true":   token.TRUE,
	"var":    token.VAR,
	"while":  token.WHILE,
}
