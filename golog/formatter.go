package golog

import (
  "fmt"
  "time"
  "path"
  "bytes"
)

type FormatContext interface {
  Severity() Severity
  FileLine() (string, int)
  Fn() string
  Msg() string
  CallerTime() time.Time
}

type Formatter interface {
  Format(FormatContext) string
}

type fmtLevel struct {
  long bool
}

func FmtLevel(long bool) fmtLevel {
  return fmtLevel{long}
}

func (self fmtLevel) Format(ctx FormatContext) string {
  if self.long {
    return ctx.Severity().String()
  }
  return ctx.Severity().Single()
}

type fmtDate struct {
  datefmt string
}

func FmtDate(datefmt string) fmtDate {
  return fmtDate{datefmt}
}

func (self fmtDate) Format(ctx FormatContext) string {
  return ctx.CallerTime().Format(self.datefmt)
}

type fmtFile struct {
  long bool
}

func FmtFile(long bool) fmtFile {
  return fmtFile{long}
}

func (self fmtFile) Format(ctx FormatContext) string {
  file, _ := ctx.FileLine()
  if self.long {
    return file
  }
  return path.Base(file)
}

type fmtLine struct {}

func FmtLine() fmtLine {
  return fmtLine{}
}

func (self fmtLine) Format(ctx FormatContext) string {
  _, line := ctx.FileLine()
  return fmt.Sprint(line)
}

type fmtFunc struct {}

func FmtFunc() fmtFunc {
  return fmtFunc{}
}

func (self fmtFunc) Format(ctx FormatContext) string {
  return fmt.Sprintf("%s()", ctx.Fn())
}

type fmtMsg struct {}

func FmtMsg() fmtMsg {
  return fmtMsg{}
}

func (self fmtMsg) Format(ctx FormatContext) string {
  return ctx.Msg()
}

type fmtString struct {
  data string
}

func FmtString(data string) fmtString {
  return fmtString{data}
}

func (self fmtString) Format(_ FormatContext) string {
  return self.data
}

type FormatSeq []Formatter

type formatSequencer struct {
  format_seq []Formatter
}

func FormatSequencer(format_seq FormatSeq) *formatSequencer {
  return &formatSequencer{format_seq: format_seq}
}

func (self *formatSequencer) Format(ctx FormatContext) string {
  var formatted bytes.Buffer

  for _, f := range self.format_seq {
    formatted.WriteString(f.Format(ctx))
  }
  return formatted.String()
}
