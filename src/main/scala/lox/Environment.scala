package lox

import scala.collection.mutable

final class Environment(enclosing: Environment = null) {

  private val values = mutable.Map.empty[String, Any]

  def define(name: String, value: Any): Unit =
    values.put(name, value)

  def get(name: Token): Any = {
    values.getOrElse(
      name.lexeme,
      if (enclosing != null)
        enclosing.get(name)
      else
        throw Interpreter.Error(name, s"Undefined variable '${name.lexeme}'."),
    )
  }

  def assign(name: Token, value: Any): Unit = {
    if (values.contains(name.lexeme)) {
      values.put(name.lexeme, value)
    } else if (enclosing != null) {
      enclosing.assign(name, value)
    } else {
      throw Interpreter.Error(name, s"Undefined variable '${name.lexeme}'.")
    }
  }

}
