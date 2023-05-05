package lox

import java.io.{BufferedReader, InputStreamReader}
import java.nio.file.{Files, Paths}

object Main {

  object ExitCode {
    val UsageError = 64
    val InvalidInput = 65
    val RuntimeError = 70
  }

  private val interpreter = new Interpreter

  private var hadError = false
  private var hadRuntimeError = false

  def main(args: Array[String]): Unit =
    if (args.length > 1) {
      println("Usage: jlox [script]");
      sys.exit(ExitCode.UsageError)
    } else if (args.length == 1) {
      run(Files.readString(Paths.get(args(0))))
      if (hadError) {
        sys.exit(ExitCode.InvalidInput)
      }
      if (hadRuntimeError) {
        sys.exit(ExitCode.RuntimeError)
      }
    } else {
      runPrompt()
    }

  private def runPrompt(): Unit = {
    val input = new InputStreamReader(System.in)
    val reader = new BufferedReader(input)

    while (true) {
      print("> ")
      val line = reader.readLine
      if (line == null) return
      run(line)
      hadError = false
    }
  }

  private def run(source: String): Unit = {
    val scanner = new Scanner(source)
    val tokens = scanner.scanTokens()
    val parser = new Parser(tokens)
    val statements = parser.parse()

    // Stop if there was a syntax error.
    if (!hadError) {
      interpreter.interpret(statements)
    }

  }

  private[lox] def error(line: Int, message: String): Unit = {
    report(line, "", message)
  }

  private[lox] def error(token: Token, message: String): Unit = {
    val location =
      if (token.tokenType == TokenType.Eof) " at end"
      else s"at '${token.lexeme}'"
    report(token.line, location, message)
  }

  def runtimeError(error: Interpreter.Error): Unit = {
    System.err.println(error.getMessage + "\n[line " + error.token.line + "]")
    hadRuntimeError = true
  }

  private def report(line: Int, where: String, message: String): Unit = {
    System.err.println("[line " + line + "] Error" + where + ": " + message)
    hadError = true
  }

}
