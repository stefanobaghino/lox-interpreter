package lox

object Parser {

  private final class Error() extends RuntimeException

}

final class Parser(tokens: List[Token]) {

  private var current = 0

  def parse(): List[Statement] = {
    val statements = List.newBuilder[Statement]
    while (!reachedEnd) statements += declaration()
    statements.result()
  }

  private def declaration(): Statement =
    try {
      if (matching(TokenType.Fun)) functionDeclaration("function")
      else if (matching(TokenType.Var)) variableDeclaration()
      else statement()
    } catch {
      case _: Parser.Error =>
        synchronize()
        null
    }

  private def functionDeclaration(kind: String): Statement = {
    val name: Token = consume(TokenType.Identifier, s"Expect $kind name.")
    consume(TokenType.LeftParen, s"Expect '(' after $kind name.")
    val parameters = List.newBuilder[Token]
    var parametersNumber = 0
    if (!check(TokenType.RightParen)) {
      do {
        if (parametersNumber >= 255) {
          error(peek(), "Can't have more than 255 parameters.")
        }
        parameters += consume(TokenType.Identifier, "Expect parameter name.")
        parametersNumber += 1
      } while (matching(TokenType.Comma))

    }
    consume(TokenType.RightParen, "Expect ')' after parameters.")
    consume(TokenType.LeftBrace, "Expect '{' before " + kind + " body.")
    Statement.Fun(name, parameters.result(), block())
  }

  private def variableDeclaration(): Statement = {
    val name: Token = consume(TokenType.Identifier, "Expect variable name.")
    val initializer: Expr =
      if (matching(TokenType.Equal)) {
        expression()
      } else {
        null
      }

    consume(TokenType.Semicolon, "Expect ';' after variable declaration.")
    Statement.Variable(name, initializer)
  }

  private def statement(): Statement = {
    if (matching(TokenType.For)) forStatement()
    else if (matching(TokenType.If)) ifStatement()
    else if (matching(TokenType.Print)) printStatement()
    else if (matching(TokenType.Return)) returnStatement()
    else if (matching(TokenType.While)) whileStatement()
    else if (matching(TokenType.LeftBrace)) Statement.Block(block())
    else expressionStatement()
  }

  private def forStatement(): Statement = {
    consume(TokenType.LeftParen, "Expect '(' after 'for'.")
    val pre =
      if (matching(TokenType.Semicolon)) null
      else if (matching(TokenType.Var)) variableDeclaration()
      else expressionStatement()
    val condition =
      if (check(TokenType.Semicolon)) Expr.Literal(true)
      else expression()
    consume(TokenType.Semicolon, "Expect ';' after 'for' condition.")
    val post =
      if (check(TokenType.RightParen)) null
      else expression()
    consume(TokenType.RightParen, "Expect ')' after 'for' clauses.")
    var body = statement()
    if (post != null)
      body = Statement.Block(List(body, Statement.Expression(post)))
    body = Statement.While(condition, body)
    if (pre != null)
      body = Statement.Block(List(pre, body))
    body
  }

  private def whileStatement(): Statement = {
    consume(TokenType.LeftParen, "Expect '(' after 'while'.")
    val condition = expression()
    consume(TokenType.RightParen, "Expect ')' after 'while' condition.")
    val body = statement()
    Statement.While(condition, body)
  }

  private def ifStatement(): Statement = {
    consume(TokenType.LeftParen, "Expect '(' after 'if'.")
    val condition = expression()
    consume(TokenType.RightParen, "Expect ')' after 'if' condition.")
    val thenBranch = statement()
    val elseBranch = if (matching(TokenType.Else)) statement() else null
    Statement.If(condition, thenBranch, elseBranch)
  }

  private def block(): List[Statement] = {
    val statements = List.newBuilder[Statement]
    while (!check(TokenType.RightBrace) && !reachedEnd()) {
      statements += declaration()
    }
    consume(TokenType.RightBrace, "Expect '}' after block.")
    statements.result()
  }

  private def printStatement(): Statement.Print = {
    val value = expression()
    consume(TokenType.Semicolon, "Expect ';' after value.")
    Statement.Print(value)
  }

  private def returnStatement(): Statement.Return = {
    val keyword: Token = previous()
    val value: Expr = if (check(TokenType.Semicolon)) null else expression()
    consume(TokenType.Semicolon, "Expect ';' after return value.")
    Statement.Return(keyword, value)
  }

  private def expressionStatement(): Statement.Expression = {
    val value = expression()
    consume(TokenType.Semicolon, "Expect ';' after value.")
    Statement.Expression(value)
  }

  private def expression(): Expr = assignment()

  private def assignment(): Expr = {
    val expr = or()

    if (matching(TokenType.Equal)) {
      val equals = previous()
      val value = assignment()
      expr match {
        case Expr.Variable(name) => return Expr.Assign(name, value)
        case _ => error(equals, "Invalid assignment target.")
      }
    }

    expr
  }

  private def or(): Expr = {
    var expr: Expr = and()

    while (matching(TokenType.Or)) {
      val operator: Token = previous()
      val right: Expr = and()
      expr = Expr.Logical(expr, operator, right)
    }

    expr
  }

  private def and(): Expr = {
    var expr: Expr = equality()

    while (matching(TokenType.And)) {
      val operator: Token = previous()
      val right: Expr = equality()
      expr = Expr.Logical(expr, operator, right)
    }

    expr
  }

  private def equality(): Expr = {
    var expr: Expr = comparison()
    while (matching(TokenType.BangEqual, TokenType.EqualEqual)) {
      val operator: Token = previous()
      val right: Expr = comparison()
      expr = Expr.Binary(expr, operator, right)
    }
    expr
  }

  private def comparison(): Expr = {
    var expr: Expr = term()
    while (
      matching(
        TokenType.Greater,
        TokenType.GreaterEqual,
        TokenType.Less,
        TokenType.LessEqual,
      )
    ) {
      val operator: Token = previous()
      val right: Expr = term()
      expr = Expr.Binary(expr, operator, right)
    }
    expr
  }

  private def term(): Expr = {
    var expr: Expr = factor()
    while (matching(TokenType.Minus, TokenType.Plus)) {
      val operator: Token = previous()
      val right: Expr = factor()
      expr = Expr.Binary(expr, operator, right)
    }
    expr
  }

  private def factor(): Expr = {
    var expr: Expr = unary()
    while (matching(TokenType.Slash, TokenType.Star)) {
      val operator: Token = previous()
      val right: Expr = unary()
      expr = Expr.Binary(expr, operator, right)
    }
    expr
  }

  private def unary(): Expr = {
    if (matching(TokenType.Bang, TokenType.Minus)) {
      val operator: Token = previous()
      val right: Expr = call()
      Expr.Unary(operator, right)
    } else {
      call()
    }
  }

  private def call(): Expr = {
    var expr = primary()

    while (matching(TokenType.LeftParen)) {
      expr = finishCall(expr)
    }

    expr
  }

  private def finishCall(callee: Expr): Expr = {
    val arguments = List.newBuilder[Expr]
    var argumentsNumber = 0
    if (!check(TokenType.RightParen)) {
      do {
        arguments += expression()
        argumentsNumber += 1
        if (argumentsNumber >= 255) {
          error(peek(), "Can't have more than 255 arguments.")
        }
      } while (matching(TokenType.Comma))
    }
    val paren = consume(TokenType.RightParen, "Expect ')' after arguments.")
    Expr.Call(callee, paren, arguments.result())
  }

  private def primary(): Expr = {
    if (matching(TokenType.True)) Expr.Literal(true)
    else if (matching(TokenType.False)) Expr.Literal(false)
    else if (matching(TokenType.Nil)) Expr.Literal(null)
    else if (matching(TokenType.Number, TokenType.String)) {
      Expr.Literal(previous().literal)
    } else if (matching(TokenType.Identifier)) {
      Expr.Variable(previous())
    } else if (matching(TokenType.LeftParen)) {
      val expr: Expr = expression()
      consume(TokenType.RightParen, "Expect ')' after expression.")
      Expr.Grouping(expr)
    } else {
      throw error(peek(), "Expect expression.")
    }
  }

  private def matching(tokenTypes: TokenType*): Boolean = {
    for (tokenType <- tokenTypes) {
      if (check(tokenType)) {
        advance()
        return true
      }
    }
    false
  }

  private def consume(tokenType: TokenType, message: String): Token = {
    if (check(tokenType)) advance() else throw error(peek(), message)
  }

  private def error(token: Token, message: String): Parser.Error = {
    Main.error(token, message)
    new Parser.Error
  }

  private def synchronize(): Unit = {
    advance()

    while (!reachedEnd()) {
      if (previous().tokenType == TokenType.Semicolon) return

      peek().tokenType match {
        case TokenType.Class => return
        case TokenType.Fun => return
        case TokenType.Var => return
        case TokenType.For => return
        case TokenType.If => return
        case TokenType.While => return
        case TokenType.Print => return
        case TokenType.Return => return
        case _ =>
      }

      advance()
    }

  }

  private def check(tokenType: TokenType): Boolean =
    !reachedEnd() && peek().tokenType == tokenType

  private def advance(): Token = {
    if (!reachedEnd()) current += 1
    previous()
  }

  private def reachedEnd(): Boolean =
    peek().tokenType == TokenType.Eof

  private def peek(): Token =
    tokens(current)

  private def previous(): Token =
    tokens(current - 1)

}
