package main

import (
  "./golog"
)

func main() {
  golog.DefLogger.SetSeverity(golog.SeverityDebug)
  golog.DefLogger.Debug("Test debug log")
  golog.DefLogger.Info("Test info log")
  golog.DefLogger.Warning("Test warning log")
  golog.DefLogger.Error("Test error log")
  golog.DefLogger.Fatal("Test fatal log")
}
