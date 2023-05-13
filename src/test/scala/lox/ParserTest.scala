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
      val statements =
        new Parser(List(primary, Tokens.Semicolon, Tokens.Eof)).parse()
      val expected = List(Statement.Expression(Expr.Literal(primary.literal)))
      assertEquals(statements, expected)
    }

  }

  test("grouping primaries") {

    val statements = new Parser(
      List(
        Tokens.LeftParen,
        Tokens.Nil,
        Tokens.RightParen,
        Tokens.Semicolon,
        Tokens.Eof,
      )
    ).parse()

    val expected = List(Statement.Expression(Expr.Grouping(Expr.Literal(null))))

    assertEquals(statements, expected)

  }

  test("equality") {

    val statements = new Parser(
      List(
        Tokens.True,
        Tokens.EqualEqual,
        Tokens.True,
        Tokens.BangEqual,
        Tokens.True,
        Tokens.Semicolon,
        Tokens.Eof,
      )
    ).parse()

    val expected =
      List(
        Statement.Expression(
          Expr.Binary(
            Expr.Binary(
              Expr.Literal(true),
              Tokens.EqualEqual,
              Expr.Literal(true),
            ),
            Tokens.BangEqual,
            Expr.Literal(true),
          )
        )
      )

    assertEquals(statements, expected)
  }

  test("comparison") {

    // true > true >= true < true <= true;
    val statements = new Parser(
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
        Tokens.Semicolon,
        Tokens.Eof,
      )
    ).parse()

    val expected =
      List(
        Statement.Expression(
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
        )
      )

    assertEquals(statements, expected)

  }

  test("comparison has precedence over equality") {

    // true != true > true;
    val statements = new Parser(
      List(
        Tokens.True,
        Tokens.BangEqual,
        Tokens.True,
        Tokens.Greater,
        Tokens.True,
        Tokens.Semicolon,
        Tokens.Eof,
      )
    ).parse()

    val expected = List(
      Statement.Expression(
        Expr.Binary(
          Expr.Literal(true),
          Tokens.BangEqual,
          Expr.Binary(Expr.Literal(true), Tokens.Greater, Expr.Literal(true)),
        )
      )
    )

    assertEquals(statements, expected)
  }

  test("unary") {

    val expr = new Parser(
      List(Tokens.Minus, Tokens.Number, Tokens.Semicolon, Tokens.Eof)
    ).parse()

    val expected =
      List(Statement.Expression(Expr.Unary(Tokens.Minus, Expr.Literal(1.0))))

    assertEquals(expr, expected)

  }

  test("unary has precedence over comparison") {

    val statements = new Parser(
      List(
        Tokens.Number,
        Tokens.BangEqual,
        Tokens.Minus,
        Tokens.Number,
        Tokens.Semicolon,
        Tokens.Eof,
      )
    ).parse()

    val expected = List(
      Statement.Expression(
        Expr.Binary(
          Expr.Literal(1.0),
          Tokens.BangEqual,
          Expr.Unary(Tokens.Minus, Expr.Literal(1.0)),
        )
      )
    )

    assertEquals(statements, expected)

  }

  test("factor/term") {

    // (1 - 1 * 1) + -1 / 1;
    val statements = new Parser(
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
        Tokens.Semicolon,
        Tokens.Eof,
      )
    ).parse()

    val expected =
      List(
        Statement.Expression(
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
        )
      )

    assertEquals(statements, expected)
  }

  test("declare/assign/block/scope/print") {
    // var a = 1; { a = "foo"; print a; } print a;
    val statements = new Parser(
      List(
        Tokens.Var,
        Tokens.Identifier,
        Tokens.Equal,
        Tokens.Number,
        Tokens.Semicolon,
        Tokens.LeftBrace,
        Tokens.Identifier,
        Tokens.Equal,
        Tokens.String,
        Tokens.Semicolon,
        Tokens.Print,
        Tokens.Identifier,
        Tokens.Semicolon,
        Tokens.RightBrace,
        Tokens.Print,
        Tokens.Identifier,
        Tokens.Semicolon,
        Tokens.Eof,
      )
    ).parse()

    val expected =
      List(
        Statement.Variable(Tokens.Identifier, Expr.Literal(1.0)),
        Statement.Block(
          List(
            Statement.Expression(
              Expr.Assign(Tokens.Identifier, Expr.Literal("foo"))
            ),
            Statement.Print(Expr.Variable(Tokens.Identifier)),
          )
        ),
        Statement.Print(Expr.Variable(Tokens.Identifier)),
      )

    assertEquals(statements, expected)
  }

  test("if") {
    // if (a) true;
    val statements = new Parser(
      List(
        Tokens.If,
        Tokens.LeftParen,
        Tokens.Identifier,
        Tokens.RightParen,
        Tokens.True,
        Tokens.Semicolon,
        Tokens.Eof,
      )
    ).parse()

    val expected =
      List(
        Statement.If(
          Expr.Variable(Tokens.Identifier),
          Statement.Expression(Expr.Literal(true)),
          null,
        )
      )

    assertEquals(statements, expected)

  }

  test("if/else") {
    // if (a) true; else false;
    val statements = new Parser(
      List(
        Tokens.If,
        Tokens.LeftParen,
        Tokens.Identifier,
        Tokens.RightParen,
        Tokens.True,
        Tokens.Semicolon,
        Tokens.Else,
        Tokens.False,
        Tokens.Semicolon,
        Tokens.Eof,
      )
    ).parse()

    val expected =
      List(
        Statement.If(
          Expr.Variable(Tokens.Identifier),
          Statement.Expression(Expr.Literal(true)),
          Statement.Expression(Expr.Literal(false)),
        )
      )

    assertEquals(statements, expected)

  }

  test("and has precedence over or") {
    // true or true and true;
    val statements = new Parser(
      List(
        Tokens.True,
        Tokens.Or,
        Tokens.True,
        Tokens.And,
        Tokens.True,
        Tokens.Semicolon,
        Tokens.Eof,
      )
    ).parse()

    val expected =
      List(
        Statement.Expression(
          Expr.Logical(
            Expr.Literal(true),
            Tokens.Or,
            Expr.Logical(Expr.Literal(true), Tokens.And, Expr.Literal(true)),
          )
        )
      )

    assertEquals(statements, expected)

  }

  test("while") {
    // while (a) true;
    val statements = new Parser(
      List(
        Tokens.While,
        Tokens.LeftParen,
        Tokens.Identifier,
        Tokens.RightParen,
        Tokens.True,
        Tokens.Semicolon,
        Tokens.Eof,
      )
    ).parse()

    val expected =
      List(
        Statement.While(
          Expr.Variable(Tokens.Identifier),
          Statement.Expression(Expr.Literal(true)),
        )
      )

    assertEquals(statements, expected)

  }

  test("for(;;)") {
    val statements = new Parser(
      List(
        Tokens.For,
        Tokens.LeftParen,
        Tokens.Semicolon,
        Tokens.Semicolon,
        Tokens.RightParen,
        Tokens.True,
        Tokens.Semicolon,
        Tokens.Eof,
      )
    ).parse()

    val expected =
      List(
        Statement.While(
          Expr.Literal(true),
          Statement.Expression(Expr.Literal(true)),
        )
      )

    assertEquals(statements, expected)

  }

  test("for(a;;)") {
    val statements = new Parser(
      List(
        Tokens.For,
        Tokens.LeftParen,
        Tokens.Identifier,
        Tokens.Semicolon,
        Tokens.Semicolon,
        Tokens.RightParen,
        Tokens.True,
        Tokens.Semicolon,
        Tokens.Eof,
      )
    ).parse()

    val expected =
      List(
        Statement.Block(
          List(
            Statement.Expression(Expr.Variable(Tokens.Identifier)),
            Statement.While(
              Expr.Literal(true),
              Statement.Expression(Expr.Literal(true)),
            ),
          )
        )
      )

    assertEquals(statements, expected)

  }

  test("for(var a = 1;;)") {
    val statements = new Parser(
      List(
        Tokens.For,
        Tokens.LeftParen,
        Tokens.Var,
        Tokens.Identifier,
        Tokens.Equal,
        Tokens.Number,
        Tokens.Semicolon,
        Tokens.Semicolon,
        Tokens.RightParen,
        Tokens.True,
        Tokens.Semicolon,
        Tokens.Eof,
      )
    ).parse()

    val expected =
      List(
        Statement.Block(
          List(
            Statement.Variable(Tokens.Identifier, Expr.Literal(1.0)),
            Statement.While(
              Expr.Literal(true),
              Statement.Expression(Expr.Literal(true)),
            ),
          )
        )
      )

    assertEquals(statements, expected)

  }

  test("for(var a = 1; a < 1;)") {
    val statements = new Parser(
      List(
        Tokens.For,
        Tokens.LeftParen,
        Tokens.Var,
        Tokens.Identifier,
        Tokens.Equal,
        Tokens.Number,
        Tokens.Semicolon,
        Tokens.Identifier,
        Tokens.Less,
        Tokens.Number,
        Tokens.Semicolon,
        Tokens.RightParen,
        Tokens.True,
        Tokens.Semicolon,
        Tokens.Eof,
      )
    ).parse()

    val expected =
      List(
        Statement.Block(
          List(
            Statement.Variable(Tokens.Identifier, Expr.Literal(1.0)),
            Statement.While(
              Expr.Binary(
                Expr.Variable(Tokens.Identifier),
                Tokens.Less,
                Expr.Literal(1.0),
              ),
              Statement.Expression(Expr.Literal(true)),
            ),
          )
        )
      )

    assertEquals(statements, expected)

  }

  test("for(var a = 1; a < 1; a = a + 1)") {
    val statements = new Parser(
      List(
        Tokens.For,
        Tokens.LeftParen,
        Tokens.Var,
        Tokens.Identifier,
        Tokens.Equal,
        Tokens.Number,
        Tokens.Semicolon,
        Tokens.Identifier,
        Tokens.Less,
        Tokens.Number,
        Tokens.Semicolon,
        Tokens.Identifier,
        Tokens.Equal,
        Tokens.Identifier,
        Tokens.Plus,
        Tokens.Number,
        Tokens.RightParen,
        Tokens.True,
        Tokens.Semicolon,
        Tokens.Eof,
      )
    ).parse()

    val expected =
      List(
        Statement.Block(
          List(
            Statement.Variable(Tokens.Identifier, Expr.Literal(1.0)),
            Statement.While(
              Expr.Binary(
                Expr.Variable(Tokens.Identifier),
                Tokens.Less,
                Expr.Literal(1.0),
              ),
              Statement.Block(
                List(
                  Statement.Expression(Expr.Literal(true)),
                  Statement.Expression(
                    Expr.Assign(
                      Tokens.Identifier,
                      Expr.Binary(
                        Expr.Variable(Tokens.Identifier),
                        Tokens.Plus,
                        Expr.Literal(1.0),
                      ),
                    )
                  ),
                )
              ),
            ),
          )
        )
      )

    assertEquals(statements, expected)

  }

  test("fun a(a) { return 1; }") {
    val statements = new Parser(
      List(
        Tokens.Fun,
        Tokens.Identifier,
        Tokens.LeftParen,
        Tokens.Identifier,
        Tokens.RightParen,
        Tokens.LeftBrace,
        Tokens.Return,
        Tokens.Number,
        Tokens.Semicolon,
        Tokens.RightBrace,
        Tokens.Eof,
      )
    ).parse()

    val expected =
      List(
        Statement.Fun(
          Tokens.Identifier,
          List(Tokens.Identifier),
          List(
            Statement.Return(Tokens.Return, Expr.Literal(1.0))
          ),
        )
      )

    assertEquals(statements, expected)
  }

}
