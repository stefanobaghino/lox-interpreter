package parser

import (
	"fmt"
	"lox"
	"lox/ast"
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

func (p *Parser) NextStatement() (stmt ast.Stmt, err error) {
	defer func() {
		if e := recover(); e != nil {
			if se, ok := e.(lox.Error); ok {
				p.sync()
				err = se
			} else {
				panic(fmt.Errorf("unexpected error during parsing: %v", e))
			}
		}
	}()
	return p.statement(), nil
}

func (p *Parser) statement() ast.Stmt {
	if p.oneOf(token.EOF) {
		return p.endStatement()
	}
	if p.oneOf(token.VAR) {
		return p.varDeclStatement()
	}
	if p.oneOf(token.ASSERT) {
		return p.assertStatement()
	}
	if p.oneOf(token.PRINT) {
		return p.printStatement()
	}
	if p.oneOf(token.LEFT_BRACE) {
		return p.blockStatement()
	}
	if p.oneOf(token.WHILE) {
		return p.whileStatement()
	}
	if p.oneOf(token.IF) {
		return p.ifStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) endStatement() ast.Stmt {
	p.pop()
	return &ast.EndStmt{}
}

func (p *Parser) varDeclStatement() ast.Stmt {
	p.pop()
	p.readToken()
	name := p.expect(token.IDENTIFIER, "expected identifier after 'var'")
	var initializer ast.Expr
	if p.oneOf(token.EQUAL) {
		p.pop()
		initializer = p.expression()
	}
	p.expect(token.SEMICOLON, "expected ';' after variable declaration")
	return &ast.VarDeclStmt{Name: name, Initializer: &initializer}
}

func (p *Parser) assertStatement() ast.Stmt {
	p.pop()
	expr := p.expression()
	p.expect(token.SEMICOLON, "expected ';' after expression")
	return &ast.AssertStmt{Expression: expr}
}

func (p *Parser) printStatement() ast.Stmt {
	p.pop()
	expr := p.expression()
	p.expect(token.SEMICOLON, "expected ';' after expression")
	return &ast.PrintStmt{Expression: expr}
}

func (p *Parser) blockStatement() ast.Stmt {
	p.pop()
	statements := []ast.Stmt{}
	for !p.oneOf(token.RIGHT_BRACE) {
		statements = append(statements, p.statement())
	}
	p.expect(token.RIGHT_BRACE, "expected '}' after block")
	return &ast.BlockStmt{Statements: statements}
}

func (p *Parser) ifStatement() ast.Stmt {
	p.pop()
	p.expect(token.LEFT_PAREN, "expected '(' after 'if'")
	condition := p.expression()
	p.expect(token.RIGHT_PAREN, "expected ')' after if condition")
	thenBranch := p.statement()
	ifStmt := &ast.IfStmt{Condition: condition, ThenBranch: &thenBranch}
	if p.oneOf(token.ELSE) {
		p.pop()
		elseBranch := p.statement()
		ifStmt.ElseBranch = &elseBranch
	}
	return ifStmt
}

func (p *Parser) whileStatement() ast.Stmt {
	p.pop()
	p.expect(token.LEFT_PAREN, "expected '(' after 'while'")
	condition := p.expression()
	p.expect(token.RIGHT_PAREN, "expected ')' after while condition")
	body := p.statement()
	return &ast.WhileStmt{Condition: condition, Body: body}
}

func (p *Parser) expressionStatement() ast.Stmt {
	expr := p.expression()
	p.expect(token.SEMICOLON, "expected ';' after expression")
	return &ast.ExprStmt{Expression: expr}
}

func (p *Parser) expression() ast.Expr {
	return p.assignment()
}

func (p *Parser) assignment() ast.Expr {
	expr := p.or()

	if p.oneOf(token.EQUAL) {
		equals := p.pop()
		value := p.assignment()

		if varExpr, ok := expr.(*ast.VarExpr); ok {
			return &ast.AssignmentExpr{Name: varExpr.Name, Value: value}
		}

		panic(&SyntaxError{equals.Line, "invalid assignment target"})
	}

	return expr
}

func (p *Parser) or() ast.Expr {
	left := p.and()

	for p.oneOf(token.OR) {
		operator := p.pop()
		right := p.and()
		left = &ast.LogicalExpr{Left: left, Operator: operator, Right: right}
	}

	return left
}

func (p *Parser) and() ast.Expr {
	left := p.equality()

	for p.oneOf(token.AND) {
		operator := p.pop()
		right := p.equality()
		left = &ast.LogicalExpr{Left: left, Operator: operator, Right: right}
	}

	return left
}

func (p *Parser) equality() ast.Expr {
	left := p.comparison()

	for p.oneOf(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.pop()
		right := p.comparison()
		left = &ast.BinaryExpr{Left: left, Operator: operator, Right: right}
	}

	return left
}

func (p *Parser) comparison() ast.Expr {
	left := p.term()

	for p.oneOf(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.pop()
		right := p.term()
		left = &ast.BinaryExpr{Left: left, Operator: operator, Right: right}
	}

	return left
}

func (p *Parser) term() ast.Expr {
	left := p.factor()

	for p.oneOf(token.MINUS, token.PLUS) {
		operator := p.pop()
		right := p.factor()
		left = &ast.BinaryExpr{Left: left, Operator: operator, Right: right}
	}

	return left
}

func (p *Parser) factor() ast.Expr {
	left := p.unary()

	for p.oneOf(token.SLASH, token.STAR) {
		operator := p.pop()
		right := p.unary()
		left = &ast.BinaryExpr{Left: left, Operator: operator, Right: right}
	}

	return left
}

func (p *Parser) unary() ast.Expr {
	if p.oneOf(token.BANG, token.MINUS) {
		operator := p.pop()
		right := p.unary()
		return &ast.UnaryExpr{Operator: operator, Right: right}
	}

	return p.primary()
}

func (p *Parser) primary() ast.Expr {
	if p.oneOf(token.IDENTIFIER) {
		name := p.pop()
		return &ast.VarExpr{Name: name}
	}
	if p.oneOf(token.FALSE) {
		p.pop()
		return &ast.LiteralExpr{Value: false}
	}
	if p.oneOf(token.TRUE) {
		p.pop()
		return &ast.LiteralExpr{Value: true}
	}
	if p.oneOf(token.NIL) {
		p.pop()
		return &ast.LiteralExpr{Value: nil}
	}
	if p.oneOf(token.NUMBER, token.STRING) {
		token := p.pop()
		return &ast.LiteralExpr{Value: token.Literal}
	}

	if p.oneOf(token.LEFT_PAREN) {
		p.pop()
		group := p.expression()
		p.expect(token.RIGHT_PAREN, "expected ')' after expression")
		return &ast.GroupingExpr{Expression: group}
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

func (p *Parser) expect(t token.Type, msg string) token.Token {
	tok := p.readToken()
	if tok.Type != t {
		if tok.Type == token.EOF {
			msg = fmt.Sprintf("%s (at end)", msg)
		} else {
			msg = fmt.Sprintf("%s (at '%s')", msg, tok.Lexeme)
		}
		panic(&SyntaxError{tok.Line, msg})
	}
	p.pop()
	return tok
}

func (p *Parser) check(t token.Type) bool {
	return p.readToken().Type == t
}

func (p *Parser) pop() token.Token {
	if len(p.tokens) == 0 {
		return p.readToken()
	}
	head, tail := p.tokens[0], p.tokens[1:]
	p.tokens = tail
	return head
}

func (p *Parser) sync() {
	defer func() {
		if e := recover(); e != nil {
			// ignore lexical errors while syncing
			if _, ok := e.(*scanner.LexicalError); !ok {
				panic(e)
			}
		}
	}()
	p.pop()
	for !syncPoint(p.readToken().Type) {
		p.pop()
	}
}

func syncPoint(t token.Type) bool {
	switch t {
	case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN, token.EOF:
		return true
	default:
		return false
	}
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
