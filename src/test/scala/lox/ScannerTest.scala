package lox

final class ScannerTest extends munit.FunSuite {

  test("empty string") {
    val tokens = new Scanner("").scanTokens()
    val expected = List(Token(TokenType.Eof, "", null, 1))
    assertEquals(tokens, expected)
  }

  test("single character tokens") {
    val tokens = new Scanner("(){},.-+;*").scanTokens()
    val expected = List(
      Token(TokenType.LeftParen, "(", null, 1),
      Token(TokenType.RightParen, ")", null, 1),
      Token(TokenType.LeftBrace, "{", null, 1),
      Token(TokenType.RightBrace, "}", null, 1),
      Token(TokenType.Comma, ",", null, 1),
      Token(TokenType.Dot, ".", null, 1),
      Token(TokenType.Minus, "-", null, 1),
      Token(TokenType.Plus, "+", null, 1),
      Token(TokenType.Semicolon, ";", null, 1),
      Token(TokenType.Star, "*", null, 1),
      Token(TokenType.Eof, "", null, 1),
    )
    assertEquals(tokens, expected)
  }

  test("single character tokens with possible lookahead") {
    val tokens = new Scanner("<<=>>=!!====/").scanTokens()
    val expected = List(
      Token(TokenType.Less, "<", null, 1),
      Token(TokenType.LessEqual, "<=", null, 1),
      Token(TokenType.Greater, ">", null, 1),
      Token(TokenType.GreaterEqual, ">=", null, 1),
      Token(TokenType.Bang, "!", null, 1),
      Token(TokenType.BangEqual, "!=", null, 1),
      Token(TokenType.EqualEqual, "==", null, 1),
      Token(TokenType.Equal, "=", null, 1),
      Token(TokenType.Slash, "/", null, 1),
      Token(TokenType.Eof, "", null, 1),
    )
    assertEquals(tokens, expected)
  }

  test("new line") {
    val tokens = new Scanner("(\n)").scanTokens()
    val expected = List(
      Token(TokenType.LeftParen, "(", null, 1),
      Token(TokenType.RightParen, ")", null, 2),
      Token(TokenType.Eof, "", null, 2),
    )
    assertEquals(tokens, expected)
  }

  test("string") {
    val tokens = new Scanner("\"hello, world\"").scanTokens()
    val expected = List(
      Token(TokenType.String, "\"hello, world\"", "hello, world", 1),
      Token(TokenType.Eof, "", null, 1),
    )
    assertEquals(tokens, expected)
  }

  test("number (without point)") {
    val tokens = new Scanner("42").scanTokens()
    val expected = List(
      Token(TokenType.Number, "42", 42.0, 1),
      Token(TokenType.Eof, "", null, 1),
    )
    assertEquals(tokens, expected)
  }

  test("number (with point)") {
    val tokens = new Scanner("42.24").scanTokens()
    val expected = List(
      Token(TokenType.Number, "42.24", 42.24, 1),
      Token(TokenType.Eof, "", null, 1),
    )
    assertEquals(tokens, expected)
  }

  test("ignore white space") {
    val tokens = new Scanner("1\r2\t3 4").scanTokens()
    val expected = List(
      Token(TokenType.Number, "1", 1.0, 1),
      Token(TokenType.Number, "2", 2.0, 1),
      Token(TokenType.Number, "3", 3.0, 1),
      Token(TokenType.Number, "4", 4.0, 1),
      Token(TokenType.Eof, "", null, 1),
    )
    assertEquals(tokens, expected)
  }

  test("ignore comments") {
    val tokens = new Scanner("// don't mind me").scanTokens()
    val expected = List(Token(TokenType.Eof, "", null, 1))
    assertEquals(tokens, expected)
  }

  test("identifiers") {
    val tokens = new Scanner("foo _bar _42 b42").scanTokens()
    val expected = List(
      Token(TokenType.Identifier, "foo", null, 1),
      Token(TokenType.Identifier, "_bar", null, 1),
      Token(TokenType.Identifier, "_42", null, 1),
      Token(TokenType.Identifier, "b42", null, 1),
      Token(TokenType.Eof, "", null, 1),
    )
    assertEquals(tokens, expected)
  }

  test("keywords") {
    val tokens =
      new Scanner(Scanner.Keywords.keysIterator.mkString(" ")).scanTokens()
    val expected =
      Scanner.Keywords.view.map { case (lexeme, tokenType) =>
        Token(tokenType, lexeme, null, 1)
      }.toList :+ Token(TokenType.Eof, "", null, 1)
    assertEquals(tokens, expected)
  }

  test("identifier containing keywords") {
    val tokens = new Scanner("forfalse truewhile classreturn").scanTokens()
    val expected = List(
      Token(TokenType.Identifier, "forfalse", null, 1),
      Token(TokenType.Identifier, "truewhile", null, 1),
      Token(TokenType.Identifier, "classreturn", null, 1),
      Token(TokenType.Eof, "", null, 1),
    )
    assertEquals(tokens, expected)
  }

}
