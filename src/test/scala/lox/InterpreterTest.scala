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

}
