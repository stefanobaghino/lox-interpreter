package lox

import (
	"fmt"
)

type Parser struct {
	scanner *Scanner
	tokens  []Token
}

func NewParser(scanner *Scanner) *Parser {
	return &Parser{scanner: scanner}
}

func (p *Parser) Parse() (expr Expr, err error) {
	defer func() {
		if e := recover(); e != nil {
			if se, ok := e.(*syntaxError); ok {
				err = se
			} else {
				panic(fmt.Errorf("unexpected error during parsing: %v", e))
			}
		}
	}()
	expr = p.expression()
	return
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.oneOf(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.pop()
		right := p.comparison()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.oneOf(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.pop()
		right := p.term()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.oneOf(MINUS, PLUS) {
		operator := p.pop()
		right := p.factor()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.oneOf(SLASH, STAR) {
		operator := p.pop()
		right := p.unary()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.oneOf(BANG, MINUS) {
		operator := p.pop()
		right := p.unary()
		return &Unary{Operator: operator, Right: right}
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.oneOf(FALSE) {
		p.pop()
		return &Literal{Value: false}
	}
	if p.oneOf(TRUE) {
		p.pop()
		return &Literal{Value: true}
	}
	if p.oneOf(NIL) {
		p.pop()
		return &Literal{Value: nil}
	}
	if p.oneOf(NUMBER, STRING) {
		token := p.pop()
		return &Literal{Value: token.Literal}
	}

	if p.oneOf(LEFT_PAREN) {
		p.pop()
		expr := p.expression()
		p.expect(RIGHT_PAREN, "expected ')' after expression")
		return &Grouping{Expression: expr}
	}

	panic(&syntaxError{p.tokens[0].Line, "expected expression"})

}

func (p *Parser) oneOf(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			return true
		}
	}

	return false
}

func (p *Parser) expect(t TokenType, msg string) {
	token := p.pop()
	if token.Type != t {
		if token.Type == EOF {
			msg = fmt.Sprintf("%s (at end)", msg)
		} else {
			msg = fmt.Sprintf("%s (at '%s')", msg, token.Lexeme)
		}
		panic(&syntaxError{token.Line, msg})
	}
}

func (p *Parser) check(t TokenType) bool {
	return p.readToken().Type == t
}

func (p *Parser) pop() Token {
	head, tail := p.tokens[0], p.tokens[1:]
	p.tokens = tail
	return head
}

func (p *Parser) readToken() Token {
	return p.readTokenAhead(0)
}

func (p *Parser) readTokenAhead(offset int) Token {
	for d := offset - len(p.tokens) + 1; d > 0; d-- {
		if token, err := p.scanner.NextToken(); err != nil {
			panic(err)
		} else {
			p.tokens = append(p.tokens, token)
		}
	}
	return p.tokens[offset]
}
