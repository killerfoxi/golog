package main

import (
  "./golog"
)

func foo() {
  golog.Info("Entered foo()")
  golog.Debug("Something is fishy!")
  golog.Info("Leaving foo()")
}

func main() {
  golog.Current.SetSeverity(golog.SeverityDebug)
  golog.Debug("Test debug log")
  golog.Info("Test info log")
  golog.Warning("Test warning log")
  golog.Error("Test error log")
  golog.Fatal("Test fatal log")
  foo()
}
