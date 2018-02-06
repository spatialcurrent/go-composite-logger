package compositelogger

import (
  "compress/gzip"
  "os"
  "path/filepath"
  "strings"
)

import (
  "github.com/mitchellh/go-homedir"
  "github.com/sirupsen/logrus"
  "github.com/pkg/errors"
)

type CompositeLogger struct {
  Logs []*logrus.Logger
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

      if logConfig.Location == "stdout" {
        log.Out = os.Stdout
      } else if logConfig.Location == "stderr" {
        log.Out = os.Stderr
      } else {

        location, err := homedir.Expand(logConfig.Location)
        if err != nil {
          return compositeLogger, errors.New("Error: Could not expand home directory for path " + logConfig.Location + ".")
        }

        f, err := os.OpenFile(location, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
        if err != nil {
          return compositeLogger, err
        }

        filename := filepath.Base(location)

        if strings.HasSuffix(filename, ".gz") {
          log.Out = gzip.NewWriter(f)
        } else {
          log.Out = f
        }

      }

      logLevel, err := logrus.ParseLevel(logConfig.Level)
      if err != nil {
        return compositeLogger, nil
      }
      log.SetLevel(logLevel)

      logs = append(logs, log)
    }
    compositeLogger.Logs = logs
    return compositeLogger, nil
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
