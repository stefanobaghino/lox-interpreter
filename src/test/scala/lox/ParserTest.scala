package lox

final class ParserTest extends munit.FunSuite {

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
