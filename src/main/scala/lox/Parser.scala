package lox

object Parser {

  private final class Error() extends RuntimeException

}

final class Parser(tokens: List[Token]) {

  private var current = 0

  def parse(): Expr =
    try expression()
    catch { case _: Parser.Error => null }

  private def expression(): Expr = equality()

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
      val right: Expr = primary()
      Expr.Unary(operator, right)
    } else {
      primary()
    }
  }

  private def primary(): Expr = {
    if (matching(TokenType.True)) Expr.Literal(true)
    else if (matching(TokenType.False)) Expr.Literal(false)
    else if (matching(TokenType.Nil)) Expr.Literal(null)
    else if (matching(TokenType.Number, TokenType.String)) {
      Expr.Literal(previous().literal)
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
