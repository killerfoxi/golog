package golog

import (
  "fmt"
  "io"
  "os"
  "runtime"
  "bytes"
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

  SetSeverity(Severity)
  GetSeverity() Severity
}

type runtimeContext struct {
  severity Severity
  msg string
  skip int

  loaded bool
  file string
  line int
  fn string
}

func newRuntimeContext(s Severity, msg string) *runtimeContext {
  return &runtimeContext{severity: s, msg: msg, file: "???", fn: "???"}
}

func (self *runtimeContext) Severity() Severity {
  return self.severity
}

func (self *runtimeContext) Msg() string {
  return self.msg
}

func (self *runtimeContext) load() {
  if pc, file, line, ok := runtime.Caller(7); ok {
    self.file = file
    self.line = line
    if fn := runtime.FuncForPC(pc); fn != nil {
      self.fn = fn.Name()
    }
  }
}

func (self *runtimeContext) FileLine() (string, int) {
  if !self.loaded {
    self.load()
  }
  return self.file, self.line
}

func (self *runtimeContext) Fn() string {
  if !self.loaded {
    self.load()
  }
  return self.fn
}

type defaultLogger struct {
  severity Severity
  formatter Formatter
  out io.Writer
}

func (self *defaultLogger) SetSeverity(s Severity) {
  self.severity = s
}

func (self *defaultLogger) GetSeverity() Severity {
  return self.severity
}

func (self *defaultLogger) output(s Severity, msg string) {
  if s > self.severity {
    return
  }

  ctx := newRuntimeContext(s, msg)
  final := bytes.NewBufferString(self.formatter.Format(ctx))
  if !bytes.HasSuffix(final.Bytes(), []byte("\n")) {
    final.WriteString("\n")
  }
  self.out.Write(final.Bytes())
}

func (self *defaultLogger) Fatalf(format string, a ...interface{}) {
  self.output(SeverityFatal, fmt.Sprintf(format, a...))
}

func (self *defaultLogger) Fatal(a ...interface{}) {
  self.output(SeverityFatal, fmt.Sprint(a...))
}

func (self *defaultLogger) Errorf(format string, a ...interface{}) {
  self.output(SeverityError, fmt.Sprintf(format, a...))
}

func (self *defaultLogger) Error(a ...interface{}) {
  self.output(SeverityError, fmt.Sprint(a...))
}

func (self *defaultLogger) Warningf(format string, a ...interface{}) {
  self.output(SeverityWarning, fmt.Sprintf(format, a...))
}

func (self *defaultLogger) Warning(a ...interface{}) {
  self.output(SeverityWarning, fmt.Sprint(a...))
}

func (self *defaultLogger) Infof(format string, a ...interface{}) {
  self.output(SeverityInfo, fmt.Sprintf(format, a...))
}

func (self *defaultLogger) Info(a ...interface{}) {
  self.output(SeverityInfo, fmt.Sprint(a...))
}

func (self *defaultLogger) Debugf(format string, a ...interface{}) {
  self.output(SeverityDebug, fmt.Sprintf(format, a...))
}

func (self *defaultLogger) Debug(a ...interface{}) {
  self.output(SeverityDebug, fmt.Sprint(a...))
}

var DefLogger Logger = &defaultLogger{
  severity: SeverityInfo,
  out: os.Stderr,
  formatter: FormatSequencer(FormatSeq{
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
  }

func Fatal(data... interface{}) {
  DefLogger.Fatal(data...)
}
