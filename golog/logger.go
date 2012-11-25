package golog

import (
  "fmt"
  "os"
)

type Severity uint8

const (
  SeverityFatal Severity = iota
  SeverityError
  SeverityWarning
  SeverityInfo
  SeverityDebug
  severityMax
)

var Severities = [severityMax]string{"Fatal", "Error", "Warning", "Info", "Debug"}

func (self Severity) String() string {
  if self > severityMax - 1 {
    return "???"
  }

  return Severities[self]
}

func (self Severity) Single() string {
  return string([]rune(self.String())[0])
}

type Logger interface {
  Fatal(...interface{})
  Fatalf(string, ...interface{})
  Error(...interface{})
  Errorf(string, ...interface{})
  Warning(...interface{})
  Warningf(string, ...interface{})
  Info(...interface{})
  Infof(string, ...interface{})
  Debug(...interface{})
  Debugf(string, ...interface{})
}

type LogWritter interface {
  Output(Severity, fmt.Stringer)

  SetSeverity(Severity)
  GetSeverity() Severity
}

type logMsg struct {
  args []interface{}
}

func newLogMsg(args []interface{}) *logMsg {
  return &logMsg{args: args}
}

func (self *logMsg) String() string {
  return fmt.Sprint(self.args...)
}

type logMsgFormatted struct {
  format string
  args []interface{}
}

func newLogMsgFormatted(format string, args []interface{}) *logMsgFormatted {
  return &logMsgFormatted{format: format, args: args}
}

func (self *logMsgFormatted) String() string {
  return fmt.Sprintf(self.format, self.args...)
}

var Current LogWritter = NewLogger(SeverityInfo,
                                   FormatSequencer(FormatSeq{
                                     FmtLevel(false),
                                     FmtDate("2006-01-02 15:04:05.999999"),
                                     FmtString(" "),
                                     FmtFile(false),
                                     FmtString("#"),
                                     FmtFunc(),
                                     FmtString(":"),
                                     FmtLine(),
                                     FmtString(": "),
                                     FmtMsg()}),
                                   os.Stderr)

func Fatalf(format string, a ...interface{}) {
  Current.Output(SeverityFatal, newLogMsgFormatted(format, a))
}

func Fatal(a ...interface{}) {
  Current.Output(SeverityFatal, newLogMsg(a))
}

func Errorf(format string, a ...interface{}) {
  Current.Output(SeverityError, newLogMsgFormatted(format, a))
}

func Error(a ...interface{}) {
  Current.Output(SeverityError, newLogMsg(a))
}

func Warningf(format string, a ...interface{}) {
  Current.Output(SeverityWarning, newLogMsgFormatted(format, a))
}

func Warning(a ...interface{}) {
  Current.Output(SeverityWarning, newLogMsg(a))
}

func Infof(format string, a ...interface{}) {
  Current.Output(SeverityInfo, newLogMsgFormatted(format, a))
}

func Info(a ...interface{}) {
  Current.Output(SeverityInfo, newLogMsg(a))
}

func Debugf(format string, a ...interface{}) {
  Current.Output(SeverityDebug, newLogMsgFormatted(format, a))
}

func Debug(a ...interface{}) {
  Current.Output(SeverityDebug, newLogMsg(a))
}
