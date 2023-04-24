package lox

final case class Token(
    tokenType: TokenType,
    lexeme: String,
    literal: Any,
    line: Int,
)
