ThisBuild / version := "0.1.0-SNAPSHOT"

ThisBuild / scalaVersion := "2.13.10"

ThisBuild / coverageEnabled := true

lazy val root = (project in file("."))
  .settings(
    name := "lox-interpreter",
    libraryDependencies ++= Seq(
      "org.scalameta" %% "munit" % "0.7.29" % Test
    ),
    Test / unmanagedSourceDirectories += (Test / sourceDirectory).value / "lox",
  )
