package lox

import java.io.{BufferedReader, InputStreamReader}
import java.nio.file.{Files, Paths}

object Main {

  object ExitCode {
    val UsageError = 64
    val InvalidInput = 65
  }

  private var hadError = false

  def main(args: Array[String]): Unit =
    if (args.length > 1) {
      println("Usage: jlox [script]");
      sys.exit(ExitCode.UsageError)
    } else if (args.length == 1) {
      run(Files.readString(Paths.get(args(0))))
      if (hadError) {
        sys.exit(ExitCode.InvalidInput)
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
    // For now, just print the tokens.
    for (token <- tokens) {
      println(token)
    }
  }

  private[lox] def error(line: Int, message: String): Unit = {
    report(line, "", message)
  }

  private def report(line: Int, where: String, message: String): Unit = {
    System.err.println("[line " + line + "] Error" + where + ": " + message)
    hadError = true
  }

}
