package lox

final class AstPrinter extends Expr.Visitor[String] {

  override def visitAssign(assign: Expr.Assign): String =
    s"${parenthesize(s"${assign.name.lexeme} =", assign.value)}"

  override def visitBinary(binary: Expr.Binary): String =
    s"${parenthesize(binary.operator.lexeme, binary.left, binary.right)}"

  override def visitGrouping(grouping: Expr.Grouping): String =
    s"${parenthesize("group", grouping.expression)}"

  override def visitLiteral(literal: Expr.Literal): String =
    if (literal == null) "nil" else literal.value.toString

  override def visitLogical(logical: Expr.Logical): String =
    s"${parenthesize(logical.operator.lexeme, logical.left, logical.right)}"

  override def visitUnary(unary: Expr.Unary): String =
    s"${parenthesize(unary.operator.lexeme, unary.right)}"

  override def visitVariableLookup(variable: Expr.Variable): String =
    s"${parenthesize(variable.name.lexeme)}"

  private def parenthesize(name: String, expressions: Expr*): String = {
    val builder = new StringBuilder
    builder.append("(").append(name)
    for (expr <- expressions) {
      builder.append(" ")
      builder.append(expr.accept(this))
    }
    builder.append(")")
    builder.toString
  }

}
