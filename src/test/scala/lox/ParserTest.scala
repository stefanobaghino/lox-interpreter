package lox

private object ParserTest {

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
  }

}

final class ParserTest extends munit.FunSuite {

  import ParserTest._

  test("non-grouping primaries") {

    val primaries = List(
      Tokens.True,
      Tokens.False,
      Tokens.String,
      Tokens.Number,
      Tokens.Nil,
    )

    for (primary <- primaries) {
      val expr = new Parser(List(primary, Tokens.Eof)).parse()
      val expected = Expr.Literal(primary.literal)
      assertEquals(expr, expected)
    }

  }

  test("grouping primaries") {

    val expr = new Parser(
      List(
        Tokens.LeftParen,
        Tokens.Nil,
        Tokens.RightParen,
        Tokens.Eof,
      )
    ).parse()

    val expected = Expr.Grouping(Expr.Literal(null))

    assertEquals(expr, expected)

  }

  test("equality") {

    val expr = new Parser(
      List(
        Tokens.True,
        Tokens.EqualEqual,
        Tokens.True,
        Tokens.BangEqual,
        Tokens.True,
        Tokens.Eof,
      )
    ).parse()

    val expected =
      Expr.Binary(
        Expr.Binary(Expr.Literal(true), Tokens.EqualEqual, Expr.Literal(true)),
        Tokens.BangEqual,
        Expr.Literal(true),
      )

    assertEquals(expr, expected)
  }

  test("comparison") {

    // true > true >= true < true <= true
    val expr = new Parser(
      List(
        Tokens.True,
        Tokens.Greater,
        Tokens.True,
        Tokens.GreaterEqual,
        Tokens.True,
        Tokens.Less,
        Tokens.True,
        Tokens.LessEqual,
        Tokens.True,
        Tokens.Eof,
      )
    ).parse()

    val expected =
      Expr.Binary(
        Expr.Binary(
          Expr.Binary(
            Expr.Binary(
              Expr.Literal(true),
              Tokens.Greater,
              Expr.Literal(true),
            ),
            Tokens.GreaterEqual,
            Expr.Literal(true),
          ),
          Tokens.Less,
          Expr.Literal(true),
        ),
        Tokens.LessEqual,
        Expr.Literal(true),
      )

    assertEquals(expr, expected)

  }

  test("comparison has precedence over equality") {

    // true != true > true
    val expr = new Parser(
      List(
        Tokens.True,
        Tokens.BangEqual,
        Tokens.True,
        Tokens.Greater,
        Tokens.True,
        Tokens.Eof,
      )
    ).parse()

    val expected = Expr.Binary(
      Expr.Literal(true),
      Tokens.BangEqual,
      Expr.Binary(Expr.Literal(true), Tokens.Greater, Expr.Literal(true)),
    )

    assertEquals(expr, expected)
  }

  test("unary") {

    val expr = new Parser(List(Tokens.Minus, Tokens.Number, Tokens.Eof)).parse()

    val expected = Expr.Unary(Tokens.Minus, Expr.Literal(1.0))

    assertEquals(expr, expected)

  }

  test("unary has precedence over comparison") {

    val expr = new Parser(
      List(
        Tokens.Number,
        Tokens.BangEqual,
        Tokens.Minus,
        Tokens.Number,
        Tokens.Eof,
      )
    ).parse()

    val expected = Expr.Binary(
      Expr.Literal(1.0),
      Tokens.BangEqual,
      Expr.Unary(Tokens.Minus, Expr.Literal(1.0)),
    )

    assertEquals(expr, expected)

  }

  test("factor/term") {

    // (1 - 1 * 1) + -1 / 1
    val expr = new Parser(
      List(
        Tokens.LeftParen,
        Tokens.Number,
        Tokens.Minus,
        Tokens.Number,
        Tokens.Star,
        Tokens.Number,
        Tokens.RightParen,
        Tokens.Plus,
        Tokens.Minus,
        Tokens.Number,
        Tokens.Slash,
        Tokens.Number,
        Tokens.Eof,
      )
    ).parse()

    val expected =
      Expr.Binary(
        Expr.Grouping(
          Expr.Binary(
            Expr.Literal(1.0),
            Tokens.Minus,
            Expr.Binary(
              Expr.Literal(1.0),
              Tokens.Star,
              Expr.Literal(1.0),
            ),
          )
        ),
        Tokens.Plus,
        Expr.Binary(
          Expr.Unary(Tokens.Minus, Expr.Literal(1.0)),
          Tokens.Slash,
          Expr.Literal(1.0),
        ),
      )

    assertEquals(expr, expected)
  }

}
