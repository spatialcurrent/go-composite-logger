package compositelogger

import (
  "compress/gzip"
  "os"
  "path/filepath"
  "strings"
  "io"
  "io/ioutil"
)

import (
  "github.com/mitchellh/go-homedir"
  "github.com/sirupsen/logrus"
  "github.com/pkg/errors"
)

type CompositeLogger struct {
  Logs []*logrus.Logger
}

func OpenOutput(location string) (io.Writer, error) {

  if location == "stdout" {
    return os.Stdout, nil
  } else if location == "stderr" {
    return os.Stderr, nil
  }

  location_expanded, err := homedir.Expand(location)
  if err != nil {
    return ioutil.Discard, errors.New("Error: Could not expand home directory for path " + location + ".")
  }

  f, err := os.OpenFile(location_expanded, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
  if err != nil {
    return ioutil.Discard, err
  }

  filename := filepath.Base(location_expanded)

  if strings.HasSuffix(filename, ".gz") {
    return gzip.NewWriter(f), nil
  }

  return f, nil
}

func NewCompositeLogger(logConfigs []*LogConfig) (*CompositeLogger, error) {
    compositeLogger := &CompositeLogger{}
    logs := make([]*logrus.Logger, 0)
    for _ , logConfig := range logConfigs {
      log := logrus.New()

      if logConfig.Format == "json" {
        log.Formatter = &logrus.JSONFormatter{}
      } else {
        log.Formatter = &logrus.TextFormatter{}
      }

      out, err := OpenOutput(logConfig.Location)
      if err != nil {
        return compositeLogger, err
      }
      log.Out = out

      logLevel, err := logrus.ParseLevel(logConfig.Level)
      if err != nil {
        return compositeLogger, err
      }
      log.SetLevel(logLevel)

      logs = append(logs, log)
    }
    compositeLogger.Logs = logs
    return compositeLogger, nil
}

func NewDefaultLogger()  (*CompositeLogger, error)  {
  logConfigs := []*LogConfig{
    &LogConfig{
      Location: "stdout",
      Level: "info",
      Format: "text",
    },
    &LogConfig{
      Location: "stderr",
      Level: "warning",
      Format: "text",
    },
  }
  return NewCompositeLogger(logConfigs)
}

func (l *CompositeLogger) Add(location string, format string, level string) error {

  log := logrus.New()

  if format == "json" {
    log.Formatter = &logrus.JSONFormatter{}
  } else {
    log.Formatter = &logrus.TextFormatter{}
  }

  out, err := OpenOutput(location)
  if err != nil {
    return err
  }
  log.Out = out

  logLevel, err := logrus.ParseLevel(level)
  if err != nil {
    return err
  }
  log.SetLevel(logLevel)

  l.Logs = append(l.Logs, log)
  return nil
}

func (l *CompositeLogger) Info(msg interface{}) {
  for _, log := range l.Logs {
    log.Info(msg)
  }
}

func (l *CompositeLogger) Warn(msg interface{}) {
  for _, log := range l.Logs {
    log.Warn(msg)
  }
}

func (l *CompositeLogger) Fatal(msg interface{}) {
  l.Logs[0].Fatal(msg)
}

func (l *CompositeLogger) InfoWithFields(msg interface{}, fields map[string]interface{}) {
  for _, log := range l.Logs {
    log.WithFields(fields).Info(msg)
  }
}

func (l *CompositeLogger) WarnWithFields(msg interface{}, fields map[string]interface{}) {
  for _, log := range l.Logs {
    log.WithFields(fields).Warn(msg)
  }
}
