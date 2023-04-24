package lox

sealed abstract class TokenType

object TokenType {

  // Single-character tokens.
  case object LeftParen extends TokenType
  case object RightParen extends TokenType
  case object LeftBrace extends TokenType
  case object RightBrace extends TokenType
  case object Comma extends TokenType
  case object Dot extends TokenType
  case object Minus extends TokenType
  case object Plus extends TokenType
  case object Semicolon extends TokenType
  case object Slash extends TokenType
  case object Star extends TokenType

  // One or two character tokens.
  case object Bang extends TokenType
  case object BangEqual extends TokenType
  case object Equal extends TokenType
  case object EqualEqual extends TokenType
  case object Greater extends TokenType
  case object GreaterEqual extends TokenType
  case object Less extends TokenType
  case object LessEqual extends TokenType

  // Literals.
  case object Identifier extends TokenType
  case object String extends TokenType
  case object Number extends TokenType

  // Keywords.
  case object And extends TokenType
  case object Class extends TokenType
  case object Else extends TokenType
  case object False extends TokenType
  case object Fun extends TokenType
  case object For extends TokenType
  case object If extends TokenType
  case object Nil extends TokenType
  case object Or extends TokenType
  case object Print extends TokenType
  case object Return extends TokenType
  case object Super extends TokenType
  case object This extends TokenType
  case object True extends TokenType
  case object Var extends TokenType
  case object While extends TokenType
  case object Eof extends TokenType

}
