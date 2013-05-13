package main

import (
  "fmt"
  "./golog"
)

func foo() {
  golog.Info("Entered foo()")
  golog.Debug("Something is fishy!")
  golog.Info("Leaving foo()")
}

func main() {
  golog.Setup(golog.FlagConfig)
  golog.Debug("Test debug log")
  golog.Info("Test info log")
  golog.Warning("Test warning log")
  golog.Error("Test error log")
  foo()
  golog.WhenDebug(func () string {
    fmt.Println("I'm visible only during debug")
    return "Test when debug"
  })
  golog.Fatal("Test fatal log")
}
