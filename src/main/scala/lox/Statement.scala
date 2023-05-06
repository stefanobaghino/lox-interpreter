package lox

sealed trait Statement {

  def accept[A](visitor: Statement.Visitor[A]): A

}

object Statement {

  trait Visitor[A] {
    def visitBlock(block: Block): A
    def visitIf(ifStatement: If): A
    def visitExpression(expr: Expression): A
    def visitPrint(print: Print): A
    def visitVariableDeclaration(variable: Variable): A
    def visitWhile(whileLoop: While): A
  }

  final case class Block(statements: List[Statement]) extends Statement {
    override def accept[A](visitor: Visitor[A]): A =
      visitor.visitBlock(this)
  }
  final case class Expression(expression: Expr) extends Statement {
    override def accept[A](visitor: Visitor[A]): A =
      visitor.visitExpression(this)
  }
  final case class If(
      condition: Expr,
      thenBranch: Statement,
      elseBranch: Statement,
  ) extends Statement {
    override def accept[A](visitor: Visitor[A]): A =
      visitor.visitIf(this)
  }
  final case class Print(expression: Expr) extends Statement {
    override def accept[A](visitor: Visitor[A]): A =
      visitor.visitPrint(this)
  }
  final case class Variable(name: Token, initializer: Expr) extends Statement {
    override def accept[A](visitor: Visitor[A]): A =
      visitor.visitVariableDeclaration(this)
  }
  final case class While(condition: Expr, body: Statement) extends Statement {
    override def accept[A](visitor: Visitor[A]): A =
      visitor.visitWhile(this)
  }

}
