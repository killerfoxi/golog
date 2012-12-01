package golog

import (
  "fmt"
  "bytes"
  "time"
  "sync"
  "os"
)

type defaultLogger struct {
  severity Severity
  formatter Formatter
  out LogDispatcher
  mu sync.Mutex
}

func NewLogger(s Severity, f Formatter, out LogDispatcher) *defaultLogger {
  return &defaultLogger{severity: s, formatter: f, out: out}
}

func (self *defaultLogger) SetSeverity(s Severity) {
  self.severity = s
}

func (self *defaultLogger) GetSeverity() Severity {
  return self.severity
}

func (self *defaultLogger) output(s Severity, msg fmt.Stringer) {
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

  if err := self.out.Dispatch(s, final.Bytes()); err != nil {
    fmt.Fprintf(os.Stderr, "Couldn't dispatch msg (%s): %s\n", final, err)
    return
  }

  if s == SeverityFatal {
    fmt.Fprintf(os.Stderr, "Abort with backtrace for debugging purpose:\n%s\n", stack(false))
    os.Exit(255)
  }
}

func (self *defaultLogger) Fatalf(format string, a ...interface{}) {
  self.output(SeverityFatal, newLogMsgFormatted(format, a))
}

func (self *defaultLogger) Fatal(a ...interface{}) {
  self.output(SeverityFatal, newLogMsg(a))
}

func (self *defaultLogger) Errorf(format string, a ...interface{}) {
  self.output(SeverityError, newLogMsgFormatted(format, a))
}

func (self *defaultLogger) Error(a ...interface{}) {
  self.output(SeverityError, newLogMsg(a))
}

func (self *defaultLogger) Warningf(format string, a ...interface{}) {
  self.output(SeverityWarning, newLogMsgFormatted(format, a))
}

func (self *defaultLogger) Warning(a ...interface{}) {
  self.output(SeverityWarning, newLogMsg(a))
}

func (self *defaultLogger) Infof(format string, a ...interface{}) {
  self.output(SeverityInfo, newLogMsgFormatted(format, a))
}

func (self *defaultLogger) Info(a ...interface{}) {
  self.output(SeverityInfo, newLogMsg(a))
}

func (self *defaultLogger) Debugf(format string, a ...interface{}) {
  self.output(SeverityDebug, newLogMsgFormatted(format, a))
}

func (self *defaultLogger) Debug(a ...interface{}) {
  self.output(SeverityDebug, newLogMsg(a))
}
