package golog

import (
  "fmt"
  "os"
  "runtime"
  "errors"
  "strings"
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

var Severities = [severityMax]string{"FATAL", "ERROR", "WARNING", "INFO", "DEBUG"}

func (self Severity) String() string {
  if self > severityMax - 1 {
    return "???"
  }

  return Severities[self]
}

func (self Severity) Single() string {
  return string([]rune(self.String())[0])
}

func (self *Severity) Set(val string) error {
  var x = strings.ToUpper(val)

  for i := SeverityFatal; i <= severityMax; i++ {
    if i.String() == x || fmt.Sprintf("%d", i) == x {
      *self = i
      return nil
    }
  }
  return errors.New("Unable to parse severity")
}

// Helper function to extract the stack
func stack(all bool) []byte {
  var buf = make([]byte, 1024)
  var size int

  for {
    size = runtime.Stack(buf, all)

    if size < cap(buf) {
      break
    } else {
      buf = make([]byte, cap(buf) * 2)
    }
  }
  return buf
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
  output(Severity, fmt.Stringer)

  SetSeverity(Severity)
  GetSeverity() Severity

  SetDispatcher(dispatcher LogDispatcher)
  GetDispatcher() LogDispatcher
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

var Current Logger = nil

func init() {
  Current = NewLogger(SeverityInfo,
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
                      NewDispatchedFile(os.Stderr))
}

func Fatalf(format string, a ...interface{}) {
  Current.output(SeverityFatal, newLogMsgFormatted(format, a))
}

func Fatal(a ...interface{}) {
  Current.output(SeverityFatal, newLogMsg(a))
}

func Errorf(format string, a ...interface{}) {
  Current.output(SeverityError, newLogMsgFormatted(format, a))
}

func Error(a ...interface{}) {
  Current.output(SeverityError, newLogMsg(a))
}

func Warningf(format string, a ...interface{}) {
  Current.output(SeverityWarning, newLogMsgFormatted(format, a))
}

func Warning(a ...interface{}) {
  Current.output(SeverityWarning, newLogMsg(a))
}

func Infof(format string, a ...interface{}) {
  Current.output(SeverityInfo, newLogMsgFormatted(format, a))
}

func Info(a ...interface{}) {
  Current.output(SeverityInfo, newLogMsg(a))
}

func Debugf(format string, a ...interface{}) {
  Current.output(SeverityDebug, newLogMsgFormatted(format, a))
}

func Debug(a ...interface{}) {
  Current.output(SeverityDebug, newLogMsg(a))
}
