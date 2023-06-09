package lox

final class InterpreterTest extends munit.FunSuite {

  test("\"hel\" + \"lo\"") {
    val interpreter = new Interpreter
    val expr = Expr.Binary(Expr.Literal("hel"), Tokens.Plus, Expr.Literal("lo"))
    val result = expr.accept(interpreter)
    val expected = "hello"
    assertEquals(result, expected)
  }

  test("1 == 1") {
    val interpreter = new Interpreter
    val expr =
      Expr.Binary(Expr.Literal(1.0), Tokens.EqualEqual, Expr.Literal(1.0))
    val result = expr.accept(interpreter)
    val expected = true
    assertEquals(result, expected)
  }

  test("1 != 1") {
    val interpreter = new Interpreter
    val expr =
      Expr.Binary(Expr.Literal(1.0), Tokens.BangEqual, Expr.Literal(1.0))
    val result = expr.accept(interpreter)
    val expected = false
    assertEquals(result, expected)
  }

  test("1 > 1") {
    val interpreter = new Interpreter
    val expr =
      Expr.Binary(Expr.Literal(1.0), Tokens.Greater, Expr.Literal(1.0))
    val result = expr.accept(interpreter)
    val expected = false
    assertEquals(result, expected)
  }

  test("1 >= 1") {
    val interpreter = new Interpreter
    val expr =
      Expr.Binary(Expr.Literal(1.0), Tokens.GreaterEqual, Expr.Literal(1.0))
    val result = expr.accept(interpreter)
    val expected = true
    assertEquals(result, expected)
  }

  test("1 < 1") {
    val interpreter = new Interpreter
    val expr =
      Expr.Binary(Expr.Literal(1.0), Tokens.Less, Expr.Literal(1.0))
    val result = expr.accept(interpreter)
    val expected = false
    assertEquals(result, expected)
  }

  test("1 <= 1") {
    val interpreter = new Interpreter
    val expr =
      Expr.Binary(Expr.Literal(1.0), Tokens.LessEqual, Expr.Literal(1.0))
    val result = expr.accept(interpreter)
    val expected = true
    assertEquals(result, expected)
  }

  test("1 > 0") {
    val interpreter = new Interpreter
    val expr = Expr.Binary(Expr.Literal(1.0), Tokens.Greater, Expr.Literal(0.0))
    val result = expr.accept(interpreter)
    val expected = true
    assertEquals(result, expected)
  }

  test("!true") {
    val interpreter = new Interpreter
    val expr = Expr.Unary(Tokens.Bang, Expr.Literal(true))
    val result = expr.accept(interpreter)
    val expected = false
    assertEquals(result, expected)
  }

  test("!nil") {
    val interpreter = new Interpreter
    val expr = Expr.Unary(Tokens.Bang, Expr.Literal(null))
    val result = expr.accept(interpreter)
    val expected = true
    assertEquals(result, expected)
  }

  test("(1 - 1 * 1) + -1 / 1") {
    val interpreter = new Interpreter
    val expr = {
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
    }
    val result = expr.accept(interpreter)
    val expected = -1.0
    assertEquals(result, expected)
  }

  test("-\"muffin\"") {
    val interpreter = new Interpreter
    val expr = Expr.Unary(Tokens.Minus, Expr.Literal("muffin"))
    val error = intercept[Interpreter.Error] {
      expr.accept(interpreter)
    }
    assertEquals(error.token, Tokens.Minus)
  }

  test("true + false") {
    val interpreter = new Interpreter
    val expr = Expr.Binary(Expr.Literal(true), Tokens.Plus, Expr.Literal(false))
    val error = intercept[Interpreter.Error] {
      expr.accept(interpreter)
    }
    assertEquals(error.token, Tokens.Plus)
  }

  test("var a = 1; { var b = 1; a = a + b; }") {
    val B = Token(TokenType.Identifier, "b", null, 1)
    val interpreter = new Interpreter
    val statements =
      List(
        Statement.Variable(Tokens.Identifier, Expr.Literal(1.0)),
        Statement.Block(
          List(
            Statement.Variable(B, Expr.Literal(1.0)),
            Statement.Expression(
              Expr.Assign(
                Tokens.Identifier,
                Expr.Binary(
                  Expr.Variable(Tokens.Identifier),
                  Tokens.Plus,
                  Expr.Variable(B),
                ),
              )
            ),
          )
        ),
      )

    interpreter.interpret(statements)

    val a =
      Expr.Variable(Tokens.Identifier).accept(interpreter).asInstanceOf[Double]

    assertEquals(a, 2.0)
  }

  test("var a = 1; if (a == 2) a = false; else a = true;") {
    val interpreter = new Interpreter
    val statements =
      List(
        Statement.Variable(Tokens.Identifier, Expr.Literal(1.0)),
        Statement.If(
          Expr.Binary(
            Expr.Literal(Tokens.Identifier),
            Tokens.EqualEqual,
            Expr.Literal(2.0),
          ),
          Statement.Expression(
            Expr.Assign(Tokens.Identifier, Expr.Literal(false))
          ),
          Statement.Expression(
            Expr.Assign(Tokens.Identifier, Expr.Literal(true))
          ),
        ),
      )

    interpreter.interpret(statements)

    val a =
      Expr.Variable(Tokens.Identifier).accept(interpreter).asInstanceOf[Boolean]

    assert(a)
  }

  test("var a = 1; if (a or a = 2) true;") {
    val interpreter = new Interpreter
    val statements =
      List(
        Statement.Variable(Tokens.Identifier, Expr.Literal(1.0)),
        Statement.If(
          Expr.Logical(
            Expr.Variable(Tokens.Identifier),
            Tokens.Or,
            Expr.Assign(Tokens.Identifier, Expr.Literal(2.0)),
          ),
          Statement.Expression(Expr.Literal(true)),
          null,
        ),
      )

    interpreter.interpret(statements)

    val a =
      Expr.Variable(Tokens.Identifier).accept(interpreter).asInstanceOf[Double]

    assertEquals(a, 1.0)
  }

  test("var a = false; if (a and a = 2) true;") {
    val interpreter = new Interpreter
    val statements =
      List(
        Statement.Variable(Tokens.Identifier, Expr.Literal(false)),
        Statement.If(
          Expr.Logical(
            Expr.Variable(Tokens.Identifier),
            Tokens.And,
            Expr.Assign(Tokens.Identifier, Expr.Literal(2.0)),
          ),
          Statement.Expression(Expr.Literal(true)),
          null,
        ),
      )

    interpreter.interpret(statements)

    val a =
      Expr.Variable(Tokens.Identifier).accept(interpreter).asInstanceOf[Boolean]

    assert(!a)
  }

  test("var a = 1; while (a <= 10) a = a + 1;") {
    val interpreter = new Interpreter
    val statements =
      List(
        Statement.Variable(Tokens.Identifier, Expr.Literal(1.0)),
        Statement.While(
          Expr.Binary(
            Expr.Variable(Tokens.Identifier),
            Tokens.LessEqual,
            Expr.Literal(10.0),
          ),
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
        ),
      )

    interpreter.interpret(statements)

    val a =
      Expr.Variable(Tokens.Identifier).accept(interpreter).asInstanceOf[Double]

    assertEquals(a, 11.0)
  }

  test("var a = 1; a();") {
    val interpreter = new Interpreter
    val statements =
      List(
        Statement.Variable(Tokens.Identifier, Expr.Literal(1.0)),
        Statement.Expression(
          Expr.Call(
            Expr.Variable(Tokens.Identifier),
            Tokens.LeftParen,
            List.empty,
          )
        ),
      )

    val error =
      intercept[Interpreter.Error] {
        for (statement <- statements) {
          statement.accept(interpreter)
        }
      }

    assertEquals(error.token, Tokens.LeftParen)
  }

  test("fun a(a) { } a();") {
    val interpreter = new Interpreter
    val statements =
      List(
        Statement.Fun(
          Tokens.Identifier,
          List(Tokens.Identifier),
          List.empty,
        ),
        Statement.Expression(
          Expr.Call(
            Expr.Variable(Tokens.Identifier),
            Tokens.LeftParen,
            List.empty,
          )
        ),
      )

    val error =
      intercept[Interpreter.Error] {
        for (statement <- statements) {
          statement.accept(interpreter)
        }
      }

    assertEquals(error.token, Tokens.LeftParen)
  }

  test("fun a() { return 1; } a = a();") {

    val interpreter = new Interpreter
    val statements =
      List(
        Statement.Fun(
          Tokens.Identifier,
          List.empty,
          List(
            Statement.Return(Tokens.Return, Expr.Literal(1.0))
          ),
        ),
        Statement.Expression(
          Expr.Assign(
            Tokens.Identifier,
            Expr.Call(
              Expr.Variable(Tokens.Identifier),
              Tokens.LeftParen,
              List.empty,
            ),
          )
        ),
      )

    interpreter.interpret(statements)

    val a =
      Expr.Variable(Tokens.Identifier).accept(interpreter).asInstanceOf[Double]

    assertEquals(a, 1.0)
  }

}
