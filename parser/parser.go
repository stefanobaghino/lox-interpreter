package parser

import (
	"fmt"
	"lox"
	"lox/expr"
	"lox/scanner"
	"lox/token"
)

type SyntaxError struct {
	line    int
	message string
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf("syntax error on line %d: %s", e.line, e.message)
}

func (e SyntaxError) Line() int {
	return e.line
}

type Parser struct {
	scanner *scanner.Scanner
	tokens  []token.Token
}

func NewParser(scanner *scanner.Scanner) *Parser {
	return &Parser{scanner: scanner}
}

func (p *Parser) Parse() (expr expr.Expr, err error) {
	defer func() {
		if e := recover(); e != nil {
			if se, ok := e.(lox.Error); ok {
				err = se
			} else {
				panic(fmt.Errorf("unexpected error during parsing: %v", e))
			}
		}
	}()
	expr = p.expression()
	return
}

func (p *Parser) expression() expr.Expr {
	return p.equality()
}

func (p *Parser) equality() expr.Expr {
	left := p.comparison()

	for p.oneOf(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.pop()
		right := p.comparison()
		left = &expr.Binary{Left: left, Operator: operator, Right: right}
	}

	return left
}

func (p *Parser) comparison() expr.Expr {
	left := p.term()

	for p.oneOf(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.pop()
		right := p.term()
		left = &expr.Binary{Left: left, Operator: operator, Right: right}
	}

	return left
}

func (p *Parser) term() expr.Expr {
	left := p.factor()

	for p.oneOf(token.MINUS, token.PLUS) {
		operator := p.pop()
		right := p.factor()
		left = &expr.Binary{Left: left, Operator: operator, Right: right}
	}

	return left
}

func (p *Parser) factor() expr.Expr {
	left := p.unary()

	for p.oneOf(token.SLASH, token.STAR) {
		operator := p.pop()
		right := p.unary()
		left = &expr.Binary{Left: left, Operator: operator, Right: right}
	}

	return left
}

func (p *Parser) unary() expr.Expr {
	if p.oneOf(token.BANG, token.MINUS) {
		operator := p.pop()
		right := p.unary()
		return &expr.Unary{Operator: operator, Right: right}
	}

	return p.primary()
}

func (p *Parser) primary() expr.Expr {
	if p.oneOf(token.FALSE) {
		p.pop()
		return &expr.Literal{Value: false}
	}
	if p.oneOf(token.TRUE) {
		p.pop()
		return &expr.Literal{Value: true}
	}
	if p.oneOf(token.NIL) {
		p.pop()
		return &expr.Literal{Value: nil}
	}
	if p.oneOf(token.NUMBER, token.STRING) {
		token := p.pop()
		return &expr.Literal{Value: token.Literal}
	}

	if p.oneOf(token.LEFT_PAREN) {
		p.pop()
		group := p.expression()
		p.expect(token.RIGHT_PAREN, "expected ')' after expression")
		return &expr.Grouping{Expression: group}
	}

	panic(&SyntaxError{p.tokens[0].Line, "expected expression"})

}

func (p *Parser) oneOf(types ...token.Type) bool {
	for _, t := range types {
		if p.check(t) {
			return true
		}
	}

	return false
}

func (p *Parser) expect(t token.Type, msg string) {
	tok := p.pop()
	if tok.Type != t {
		if tok.Type == token.EOF {
			msg = fmt.Sprintf("%s (at end)", msg)
		} else {
			msg = fmt.Sprintf("%s (at '%s')", msg, tok.Lexeme)
		}
		panic(&SyntaxError{tok.Line, msg})
	}
}

func (p *Parser) check(t token.Type) bool {
	return p.readToken().Type == t
}

func (p *Parser) pop() token.Token {
	head, tail := p.tokens[0], p.tokens[1:]
	p.tokens = tail
	return head
}

func (p *Parser) readToken() token.Token {
	return p.readTokenAhead(0)
}

func (p *Parser) readTokenAhead(offset int) token.Token {
	for d := offset - len(p.tokens) + 1; d > 0; d-- {
		if token, err := p.scanner.NextToken(); err != nil {
			panic(err)
		} else {
			p.tokens = append(p.tokens, token)
		}
	}
	return p.tokens[offset]
}
