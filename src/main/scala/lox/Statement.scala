package lox

sealed trait Statement {

  def accept[A](visitor: Statement.Visitor[A]): A

}

object Statement {

  trait Visitor[A] {
    def visitBlock(block: Block): A
    def visitExpression(expr: Expression): A
    def visitPrint(print: Print): A
    def visitVariableDeclaration(variable: Variable): A
  }

  final case class Block(statements: List[Statement]) extends Statement {
    override def accept[A](visitor: Visitor[A]): A =
      visitor.visitBlock(this)
  }
  final case class Expression(expression: Expr) extends Statement {
    override def accept[A](visitor: Visitor[A]): A =
      visitor.visitExpression(this)
  }
  final case class Print(expression: Expr) extends Statement {
    override def accept[A](visitor: Visitor[A]): A =
      visitor.visitPrint(this)
  }
  final case class Variable(name: Token, initializer: Expr) extends Statement {
    override def accept[A](visitor: Visitor[A]): A =
      visitor.visitVariableDeclaration(this)
  }

}
