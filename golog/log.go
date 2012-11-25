package golog

import (
  "fmt"
  "io"
  "os"
  "runtime"
  "bytes"
  "time"
  "sync"
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

type runtimeContext struct {
  severity Severity
  msg string
  skip int

  file string
  line int
  fn string
  caller_time time.Time
}

func newRuntimeContext(s Severity, msg string, caller_time time.Time) *runtimeContext {
  ctx := &runtimeContext{severity: s, msg: msg, file: "???", fn: "???",
                         caller_time: caller_time, skip: 4}
  ctx.load()
  return ctx
}

func (self *runtimeContext) Severity() Severity {
  return self.severity
}

func (self *runtimeContext) Msg() string {
  return self.msg
}

func (self *runtimeContext) load() {
  if pc, file, line, ok := runtime.Caller(self.skip); ok {
    self.file = file
    self.line = line
    if fn := runtime.FuncForPC(pc); fn != nil {
      self.fn = fn.Name()
    }
  }
}

func (self *runtimeContext) FileLine() (string, int) {
  return self.file, self.line
}

func (self *runtimeContext) Fn() string {
  return self.fn
}

func (self *runtimeContext) CallerTime() time.Time {
  return self.caller_time
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

type defaultLogger struct {
  severity Severity
  formatter Formatter
  out io.Writer
  mu sync.Mutex
}

func NewLogger(s Severity, f Formatter, out io.Writer) *defaultLogger {
  return &defaultLogger{severity: s, formatter: f, out: out}
}

func (self *defaultLogger) SetSeverity(s Severity) {
  self.severity = s
}

func (self *defaultLogger) GetSeverity() Severity {
  return self.severity
}

func (self *defaultLogger) Output(s Severity, msg fmt.Stringer) {
  if s > self.severity {
    return
  }

  self.mu.Lock()
  defer self.mu.Unlock()

  ctx := newRuntimeContext(s, msg.String(), time.Now())
  final := bytes.NewBufferString(self.formatter.Format(ctx))
  if !bytes.HasSuffix(final.Bytes(), []byte("\n")) {
    final.WriteString("\n")
  }
  self.out.Write(final.Bytes())
}

func (self *defaultLogger) Fatalf(format string, a ...interface{}) {
  self.Output(SeverityFatal, newLogMsgFormatted(format, a))
}

func (self *defaultLogger) Fatal(a ...interface{}) {
  self.Output(SeverityFatal, newLogMsg(a))
}

func (self *defaultLogger) Errorf(format string, a ...interface{}) {
  self.Output(SeverityError, newLogMsgFormatted(format, a))
}

func (self *defaultLogger) Error(a ...interface{}) {
  self.Output(SeverityError, newLogMsg(a))
}

func (self *defaultLogger) Warningf(format string, a ...interface{}) {
  self.Output(SeverityWarning, newLogMsgFormatted(format, a))
}

func (self *defaultLogger) Warning(a ...interface{}) {
  self.Output(SeverityWarning, newLogMsg(a))
}

func (self *defaultLogger) Infof(format string, a ...interface{}) {
  self.Output(SeverityInfo, newLogMsgFormatted(format, a))
}

func (self *defaultLogger) Info(a ...interface{}) {
  self.Output(SeverityInfo, newLogMsg(a))
}

func (self *defaultLogger) Debugf(format string, a ...interface{}) {
  self.Output(SeverityDebug, newLogMsgFormatted(format, a))
}

func (self *defaultLogger) Debug(a ...interface{}) {
  self.Output(SeverityDebug, newLogMsg(a))
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
