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
  "errors"
)

type LogDispatcher interface {
  Dispatch(Severity, []byte) error
}

type dispatchedFile struct {
  fh io.Writer
}

func NewDispatchedFile(fh io.Writer) *dispatchedFile {
  return &dispatchedFile{fh: fh}
}

func (self *dispatchedFile) Dispatch(_ Severity, data []byte) error {
  _, err := self.fh.Write(data)
  return err
}

type ioMultiDispatcher struct {
  out [severityMax]io.Writer
  pushed_upwards bool
}

func NewIoMultiDispatcher(out [severityMax]io.Writer, pushed_upwards bool) *ioMultiDispatcher {
  return &ioMultiDispatcher{out: out, pushed_upwards: pushed_upwards}
}

func (self *ioMultiDispatcher) Dispatch(s Severity, data []byte) error {
  if s >= severityMax {
    return errors.New("Unknown severity")
  }

  if s == SeverityFatal {
    // Write all stack traces
    bf := bytes.NewBuffer(data)
    bf.Write(stack(true))
    data = bf.Bytes()
  }

  for i := s; i < severityMax; i++ {
    if _, err := self.out[i].Write(data); err != nil {
      return err
    }

    if !self.pushed_upwards {
      break
    }
  }
  return nil
}

var hostname = "???"
var username = "???"
var pid = os.Getpid()

func init() {
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

  fh *os.File
}

func NewLogFile(logdir string, tag string) *logFile {
  //file := &logFile{logdir: logdir, tag: tag}
  file := new(logFile)
  file.logdir = logdir
  file.tag = tag
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

func (self *logFile) Dispatch(_ Severity, data []byte) error {
  _, err := self.fh.Write(data)
  return err
}
