package golog

import (
  "fmt"
  "os"
  "runtime"
  "flag"
  "errors"
  "strings"
  "io"
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

var Current LogWritter = nil

func init() {
  var logtostderr, logalsotostderr bool
  var logdir, logfiletag string
  var logthreshold = SeverityInfo

  var output io.Writer

  fs := flag.NewFlagSet("golog", flag.ContinueOnError)
  fs.BoolVar(&logtostderr, "logtostderr", true, "Log to stderr")
  fs.BoolVar(&logalsotostderr, "logalsotostderr", false, "Log also to stderr")
  fs.Var(&logthreshold, "logthreshold", "The log threshold")
  fs.StringVar(&logdir, "logdir", ".", "Specifies the logdir")
  fs.StringVar(&logfiletag, "logfiletag", "ALL", "Specifies the logfile tag")
  if err := fs.Parse(os.Args[1:]); err != nil {
    panic(err)
  }

  if logtostderr {
    output = os.Stderr
  } else {
    multiplexerInit()
    filelog := NewIoFile(logdir, logfiletag)
    if err := filelog.Setup(); err != nil {
      panic(err)
    }

    if logalsotostderr {
      output = io.MultiWriter(filelog, os.Stderr)
    } else {
      output = filelog
    }
  }

  Current = NewLogger(logthreshold,
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
                      output)
}

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
