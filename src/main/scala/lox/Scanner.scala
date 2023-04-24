package lox

import scala.collection.mutable.ListBuffer

object Scanner {
  private[lox] val Keywords: Map[String, TokenType] = Map(
    "and" -> TokenType.And,
    "class" -> TokenType.Class,
    "else" -> TokenType.Else,
    "false" -> TokenType.False,
    "for" -> TokenType.For,
    "fun" -> TokenType.Fun,
    "if" -> TokenType.If,
    "nil" -> TokenType.Nil,
    "or" -> TokenType.Or,
    "print" -> TokenType.Print,
    "return" -> TokenType.Return,
    "super" -> TokenType.Super,
    "this" -> TokenType.This,
    "true" -> TokenType.True,
    "var" -> TokenType.Var,
    "while" -> TokenType.While,
  )
}

final class Scanner(private val source: String) {
  final private val tokens = new ListBuffer[Token]
  private var start = 0
  private var current = 0
  private var line = 1

  private def reachedEnd(): Boolean = current >= source.length

  def scanTokens(): List[Token] = {
    while (!reachedEnd()) {
      // We are at the beginning of the next lexeme.
      start = current
      scanToken()
    }
    tokens += Token(TokenType.Eof, "", null, line)
    tokens.result()
  }

  private def scanToken(): Unit =
    advance() match {
      case '(' => addToken(TokenType.LeftParen)
      case ')' => addToken(TokenType.RightParen)
      case '{' => addToken(TokenType.LeftBrace)
      case '}' => addToken(TokenType.RightBrace)
      case ',' => addToken(TokenType.Comma)
      case '.' => addToken(TokenType.Dot)
      case '-' => addToken(TokenType.Minus)
      case '+' => addToken(TokenType.Plus)
      case ';' => addToken(TokenType.Semicolon)
      case '*' => addToken(TokenType.Star)
      case '!' if followedBy('=') => addToken(TokenType.BangEqual)
      case '!' => addToken(TokenType.Bang)
      case '=' if followedBy('=') => addToken(TokenType.EqualEqual)
      case '=' => addToken(TokenType.Equal)
      case '>' if followedBy('=') => addToken(TokenType.GreaterEqual)
      case '>' => addToken(TokenType.Greater)
      case '<' if followedBy('=') => addToken(TokenType.LessEqual)
      case '<' => addToken(TokenType.Less)
      case '/' if followedBy('/') => skipUntil('\n')
      case '/' => addToken(TokenType.Slash)
      case ' ' | '\r' | '\t' => // ignore whitespace
      case '\n' => line += 1
      case '"' => string()
      case n if isDigit(n) => number()
      case a if isAlpha(a) => identifier()
      case _ => Main.error(line, "Unexpected character.")
    }

  private def identifier(): Unit = {
    while (isAlphaNumeric(peek())) advance()
    val text = source.substring(start, current)
    addToken(Scanner.Keywords.getOrElse(text, TokenType.Identifier))
  }

  private def string(): Unit = {
    while ((peek() != '"') && !reachedEnd()) {
      if (peek() == '\n') line += 1
      advance()
    }
    if (reachedEnd()) {
      Main.error(line, "Unterminated string.")
      return
    }
    advance() // The closing quote
    val value = source.substring(start + 1, current - 1) // Trim the quotes
    addToken(TokenType.String, value)
  }

  private def number(): Unit = {
    while (isDigit(peek())) advance()
    // Look for a fractional part.
    if ((peek() == '.') && isDigit(peekNext())) {
      advance() // Consume the "."
      while (isDigit(peek())) advance()
    }
    addToken(TokenType.Number, source.substring(start, current).toDouble)
  }

  private def followedBy(expected: Char): Boolean = {
    if (reachedEnd()) return false
    if (source.charAt(current) != expected) return false
    current += 1
    true
  }

  private def skipUntil(char: Char): Unit =
    while (peek() != char && !reachedEnd()) advance()

  private def peek(): Char = {
    if (reachedEnd()) return '\u0000'
    source.charAt(current)
  }

  private def peekNext(): Char = {
    if (current + 1 >= source.length) return '\u0000'
    source.charAt(current + 1)
  }

  private def isAlpha(c: Char) =
    (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'

  private def isAlphaNumeric(c: Char) = isAlpha(c) || isDigit(c)

  private def isDigit(c: Char) = c >= '0' && c <= '9'

  private def advance() = {
    val next = source.charAt(current)
    current += 1
    next
  }

  private def addToken(tokenType: TokenType): Unit =
    addToken(tokenType, null)

  private def addToken(tokenType: TokenType, literal: Any): Unit = {
    val text = source.substring(start, current)
    tokens += Token(tokenType, text, literal, line)
  }

}
