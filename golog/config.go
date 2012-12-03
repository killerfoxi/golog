package golog

import (
  "io"
  "os"
  "flag"
)

type LogConfig interface {
  NeedsSetup() bool
  Setup()

  LogToStderr() bool
  LogAlsoToStderr() bool
  LogToSingle() bool
  LogDir() string
  LogFileTag() string
  LogThreshold() Severity
}

type flagConfig struct {
  logtostderr bool
  logalsotostderr bool
  logtosingle bool
  logdir string
  logfiletag string
  logthreshold Severity
}

var FlagConfig = &flagConfig{logthreshold: SeverityInfo}

func (self *flagConfig) NeedsSetup() bool {
  return !flag.Parsed()
}

func (self *flagConfig) Setup() {
  flag.Parse()
}

func (self *flagConfig) LogToStderr() bool {
  return self.logtostderr
}

func (self *flagConfig) LogAlsoToStderr() bool {
  return self.logalsotostderr
}

func (self *flagConfig) LogToSingle() bool {
  return self.logtosingle
}

func (self *flagConfig) LogDir() string {
  return self.logdir
}

func (self *flagConfig) LogFileTag() string {
  return self.logfiletag
}

func (self *flagConfig) LogThreshold() Severity {
  return self.logthreshold
}

func init() {
  flag.BoolVar(&FlagConfig.logtostderr, "logtostderr", true, "Log to stderr")
  flag.BoolVar(&FlagConfig.logalsotostderr, "logalsotostderr", false, "Log also to stderr")
  flag.BoolVar(&FlagConfig.logtosingle, "logtosingle", false, "Log to single file")
  flag.Var(&FlagConfig.logthreshold, "logthreshold", "The log threshold")
  flag.StringVar(&FlagConfig.logdir, "logdir", ".", "Specifies the logdir")
  flag.StringVar(&FlagConfig.logfiletag, "logfiletag", "ALL", "Specifies the logfile tag (if single)")
}

func Setup(cfg LogConfig) {
  var dispatcher LogDispatcher

  if cfg.NeedsSetup() {
    cfg.Setup()
  }

  if cfg.LogToStderr() {
    dispatcher = NewDispatchedFile(os.Stderr)
  } else {
    if cfg.LogToSingle() {
      filelog := NewLogFile(cfg.LogDir(), cfg.LogFileTag())
      if err := filelog.Setup(); err != nil {
        Fatal(err)
      }

      if cfg.LogAlsoToStderr() {
        dispatcher = NewDispatchedFile(io.MultiWriter(filelog, os.Stderr))
      } else {
        dispatcher = filelog
      }
    } else {
      var filelogs [severityMax]io.Writer
      for s := SeverityFatal; s < severityMax; s++ {
        filelogs[s] = NewLogFile(cfg.LogDir(), s.String())
      }

      if cfg.LogAlsoToStderr() {
        filelogs[SeverityDebug] = io.MultiWriter(filelogs[SeverityDebug], os.Stderr)
      }
      dispatcher = NewIoMultiDispatcher(filelogs, true)
    }
  }

  Current.SetDispatcher(dispatcher)
  Current.SetSeverity(cfg.LogThreshold())
}
