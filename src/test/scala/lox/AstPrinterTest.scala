package lox

final class AstPrinterTest extends munit.FunSuite {

  test("complex expression") {
    val printer = new AstPrinter
    val expression = Expr.Binary(
      Expr.Unary(Token(TokenType.Minus, "-", null, 1), Expr.Literal(123)),
      Token(TokenType.Star, "*", null, 1),
      Expr.Grouping(Expr.Literal(45.67)),
    )
    val result = expression.accept(printer)
    val expected = "(* (- 123) (group 45.67))"
    assertEquals(result, expected)
  }

}
