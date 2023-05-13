package lox

sealed trait Callable {

  def arity: Int

  def call(interpreter: Interpreter, arguments: List[Any]): Any

}

object Callable {

  object Natives {

    object Clock extends Callable {
      override def arity: Int = 0

      override def call(interpreter: Interpreter, arguments: List[Any]): Any =
        System.currentTimeMillis.asInstanceOf[Double] / 1000

      override def toString: String = "<native fn>"
    }

  }

  final class Fun(declaration: Statement.Fun, closure: Environment)
      extends Callable {
    override val arity: Int = declaration.params.length

    override def call(interpreter: Interpreter, arguments: List[Any]): Any = {
      val env = new Environment(closure)
      val bindings = declaration.params.view.map(_.lexeme).zip(arguments)
      for ((name, value) <- bindings) {
        env.define(name, value)
      }
      try {
        interpreter.executeBlock(declaration.body, env)
      } catch {
        case Interpreter.Return(value) => return value
      }
      null
    }

    override def toString: String = s"<fn ${declaration.name.lexeme}>"
  }

}
