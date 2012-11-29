package golog

import (
  "fmt"
  "os"
  "os/user"
  "bytes"
  "io"
  "strings"
  "time"
  "path"
)

type ioMultiplexer struct {
  parent LogWritter
  out [severityMax]io.Writer
  pushed_upwards bool
}

func NewIoMultiplexer(parent LogWritter, out [severityMax]io.Writer, pushed_upwards bool) *ioMultiplexer {
  return &ioMultiplexer{parent: parent, out: out, pushed_upwards: pushed_upwards}
}

func (self *ioMultiplexer) Output(s Severity, msg fmt.Stringer) {
  if s >= severityMax {
    return
  }

  var log []byte
  if s != SeverityFatal {
    log = []byte(msg.String())
  } else {
    // Write all stack traces
    bf := bytes.NewBufferString(msg.String())
    bf.Write(stack(true))
    log = bf.Bytes()
  }

  for i := s; i < severityMax; i++ {
    self.out[i].Write(log)

    if !self.pushed_upwards {
      break
    }
  }
}

func (self *ioMultiplexer) SetSeverity(s Severity) {
  self.parent.SetSeverity(s)
}

func (self *ioMultiplexer) GetSeverity() Severity {
  return self.parent.GetSeverity()
}

var hostname = "???"
var username = "???"
var pid = os.Getpid()

func multiplexerInit() {
  if hn, err := os.Hostname(); err == nil {
    hostname = strings.SplitN(hn, ".", 2)[0]
  }
  if un, err := user.Current(); err == nil {
    username = un.Username
  }
}

type logFile struct {
  logdir string
  tag string

  // Calculated
  fh *os.File
}

func NewLogFile(logdir string, tag string) *logFile {
  file := &logFile{logdir: logdir, tag: tag}
  if err := file.Setup(); err != nil {
    panic(err)
  }
  return file
}

func (self *logFile) Setup() (err error) {
  if err = os.MkdirAll(self.logdir, 0700); err != nil {
    return
  }

  t := time.Now()
  logpath := path.Join(self.logdir, fmt.Sprintf("%s.%s.%s.log.%s.%04d%02d%02d-%02d%02d%02d.%d", path.Base(os.Args[0]), hostname, username, self.tag, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), pid))

  if self.fh, err = os.Create(logpath); err != nil {
    return
  }
  return nil
}

func (self *logFile) Write(p []byte) (int, error) {
  return self.fh.Write(p)
}
