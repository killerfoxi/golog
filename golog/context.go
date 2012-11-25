package golog

import (
  "time"
  "runtime"
)

type FormatContext interface {
  Severity() Severity
  FileLine() (string, int)
  Fn() string
  Msg() string
  CallerTime() time.Time
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

