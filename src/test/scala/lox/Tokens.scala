package lox

object Tokens {

  val Eof: Token = Token(TokenType.Eof, "", null, 1)
  val True: Token = Token(TokenType.True, "true", true, 1)
  val False: Token = Token(TokenType.False, "false", false, 1)
  val String: Token = Token(TokenType.String, "\"foo\"", "foo", 1)
  val Number: Token = Token(TokenType.Number, "1", 1.0, 1)
  val Nil: Token = Token(TokenType.Nil, "nil", null, 1)
  val LeftParen: Token = Token(TokenType.LeftParen, "(", null, 1)
  val RightParen: Token = Token(TokenType.RightParen, ")", null, 1)
  val EqualEqual: Token = Token(TokenType.EqualEqual, "==", null, 1)
  val BangEqual: Token = Token(TokenType.BangEqual, "!=", null, 1)
  val Greater: Token = Token(TokenType.Greater, ">", null, 1)
  val GreaterEqual: Token = Token(TokenType.GreaterEqual, ">=", null, 1)
  val Less: Token = Token(TokenType.Less, "<", null, 1)
  val LessEqual: Token = Token(TokenType.LessEqual, "<=", null, 1)
  val Plus: Token = Token(TokenType.Plus, "+", null, 1)
  val Minus: Token = Token(TokenType.Minus, "-", null, 1)
  val Star: Token = Token(TokenType.Star, "*", null, 1)
  val Slash: Token = Token(TokenType.Slash, "/", null, 1)
  val Bang: Token = Token(TokenType.Bang, "!", null, 1)
  val Semicolon: Token = Token(TokenType.Semicolon, ";", null, 1)
  val Var: Token = Token(TokenType.Var, "var", null, 1)
  val Identifier: Token = Token(TokenType.Identifier, "a", null, 1)
  val Equal: Token = Token(TokenType.Equal, "=", null, 1)
  val LeftBrace: Token = Token(TokenType.LeftBrace, "{", null, 1)
  val RightBrace: Token = Token(TokenType.RightBrace, "}", null, 1)
  val Print: Token = Token(TokenType.Print, "print", null, 1)
  var And: Token = Token(TokenType.And, "and", null, 1)
  var Or: Token = Token(TokenType.Or, "or", null, 1)
  var If: Token = Token(TokenType.If, "if", null, 1)
  var Else: Token = Token(TokenType.Else, "else", null, 1)
  val While: Token = Token(TokenType.While, "while", null, 1)
  val For: Token = Token(TokenType.For, "for", null, 1)
  val Fun: Token = Token(TokenType.Fun, "fun", null, 1)
  val Return: Token = Token(TokenType.Return, "return", null, 1)

}
