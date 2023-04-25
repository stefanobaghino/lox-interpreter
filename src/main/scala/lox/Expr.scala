package lox

import lox.Expr.Visitor

sealed trait Expr {

  def accept[A](visitor: Visitor[A]): A

}

object Expr {

  trait Visitor[A] {
    def visitBinary(binary: Binary): A
    def visitGrouping(grouping: Grouping): A
    def visitLiteral(literal: Literal): A
    def visitUnary(unary: Unary): A
  }

  final case class Binary(left: Expr, operator: Token, right: Expr)
      extends Expr {
    override def accept[A](visitor: Visitor[A]): A = visitor.visitBinary(this)
  }

  final case class Grouping(expression: Expr) extends Expr {
    override def accept[A](visitor: Visitor[A]): A = visitor.visitGrouping(this)
  }

  final case class Literal(value: Any) extends Expr {
    override def accept[A](visitor: Visitor[A]): A = visitor.visitLiteral(this)
  }

  final case class Unary(operator: Token, right: Expr) extends Expr {
    override def accept[A](visitor: Visitor[A]): A = visitor.visitUnary(this)
  }

}
