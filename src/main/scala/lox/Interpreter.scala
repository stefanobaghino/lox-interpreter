package lox

object Interpreter {

  final case class Error(token: Token, message: String)
      extends RuntimeException(message)

}

final class Interpreter extends Expr.Visitor[Any] {

  def interpret(expression: Expr): Unit = {
    try {
      val value = evaluate(expression)
      println(stringify(value))
    } catch {
      case error: Interpreter.Error =>
        Main.runtimeError(error)
    }
  }

  override def visitBinary(binary: Expr.Binary): Any = {
    val left = evaluate(binary.left)
    val right = evaluate(binary.right)

    binary.operator.tokenType match {
      case TokenType.Minus =>
        checkNumberOperands(binary.operator, left, right)
        left.asInstanceOf[Double] - right.asInstanceOf[Double]
      case TokenType.Slash =>
        checkNumberOperands(binary.operator, left, right)
        left.asInstanceOf[Double] / right.asInstanceOf[Double]
      case TokenType.Star =>
        checkNumberOperands(binary.operator, left, right)
        left.asInstanceOf[Double] * right.asInstanceOf[Double]
      case TokenType.Plus =>
        (left, right) match {
          case (n: Double, m: Double) => n + m
          case (s: String, t: String) => s + t
          case _ =>
            throw Interpreter.Error(
              binary.operator,
              "Operands must be two numbers or two strings.",
            );
        }
      case TokenType.Greater =>
        checkNumberOperands(binary.operator, left, right)
        left.asInstanceOf[Double] > right.asInstanceOf[Double]
      case TokenType.GreaterEqual =>
        checkNumberOperands(binary.operator, left, right)
        left.asInstanceOf[Double] >= right.asInstanceOf[Double]
      case TokenType.Less =>
        checkNumberOperands(binary.operator, left, right)
        left.asInstanceOf[Double] < right.asInstanceOf[Double]
      case TokenType.LessEqual =>
        checkNumberOperands(binary.operator, left, right)
        left.asInstanceOf[Double] <= right.asInstanceOf[Double]
      case TokenType.EqualEqual =>
        isEqual(left, right)
      case TokenType.BangEqual =>
        !isEqual(left, right)
      case _ =>
        null // unreachable
    }

  }

  override def visitGrouping(grouping: Expr.Grouping): Any =
    evaluate(grouping.expression)

  override def visitLiteral(literal: Expr.Literal): Any =
    literal.value

  override def visitUnary(unary: Expr.Unary): Any = {
    val right = evaluate(unary.right)

    unary.operator.tokenType match {
      case TokenType.Minus =>
        checkNumberOperand(unary.operator, right)
        -right.asInstanceOf[Double]
      case TokenType.Bang =>
        !isTruthy(right)
      case _ =>
        null // Unreachable
    }
  }

  private def checkNumberOperand(operator: Token, operand: Any): Unit =
    if (!operand.isInstanceOf[Double])
      throw Interpreter.Error(operator, "Operand must be a number.")

  private def checkNumberOperands(
      operator: Token,
      left: Any,
      right: Any,
  ): Unit =
    if (!left.isInstanceOf[Double] && right.isInstanceOf[Double])
      throw Interpreter.Error(operator, "Operands must be numbers.")

  private def isTruthy(stuff: Any) =
    stuff match {
      case boolean: Boolean => boolean
      case _ => stuff != null
    }

  private def isEqual(left: Any, right: Any): Boolean = {
    if (left == null && right == null) true
    else if (left == null) false
    else left == right
  }

  private def evaluate(expr: Expr) =
    expr.accept(this)

  private def stringify(stuff: Any): String = {
    if (stuff == null) return "nil"
    if (stuff.isInstanceOf[Double]) {
      var text = stuff.toString
      if (text.endsWith(".0")) text = text.substring(0, text.length - 2)
      return text
    }
    stuff.toString
  }

}
