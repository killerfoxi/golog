package golog

import (
  "fmt"
  "io"
  "bytes"
  "time"
  "sync"
)

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
