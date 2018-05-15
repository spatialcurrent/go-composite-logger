package compositelogger

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

import (
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// CompositeLogger is a struct containing an array of loggers.
type CompositeLogger struct {
	Logs []*logrus.Logger // array of loggers
}

// OpenOutput returns an io.Writer for a given location to log.
// If there is an error, returns the error.
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

// NewCompositeLogger returns a composite logger for an array of log configurations.
// If there is an error, returns the error.
func NewCompositeLogger(logConfigs []*LogConfig) (*CompositeLogger, error) {
	compositeLogger := &CompositeLogger{}
	logs := make([]*logrus.Logger, 0)
	for _, logConfig := range logConfigs {
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

// NewDefaultLogger returns a default composite logger, with stdout/info and stderr/warning.
// If there is an error, returns the error.
func NewDefaultLogger() (*CompositeLogger, error) {
	logConfigs := []*LogConfig{
		&LogConfig{
			Location: "stdout",
			Level:    "info",
			Format:   "text",
		},
		&LogConfig{
			Location: "stderr",
			Level:    "warning",
			Format:   "text",
		},
	}
	return NewCompositeLogger(logConfigs)
}

// Add adds a new logger to the composite logger, specified by the parameters.
// Returns an error, if any.
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

// Info logs a message with level INFO to all loggers.
func (l *CompositeLogger) Info(msg interface{}) {
	for _, log := range l.Logs {
		log.Info(msg)
	}
}

// Info logs a message with level WARN to all loggers.
func (l *CompositeLogger) Warn(msg interface{}) {
	for _, log := range l.Logs {
		log.Warn(msg)
	}
}

// Info logs a message with level FATAL to all loggers.
func (l *CompositeLogger) Fatal(msg interface{}) {
	l.Logs[0].Fatal(msg)
}

// InfoWithFields logs a message with level INFO and a map of fields to all loggers.
func (l *CompositeLogger) InfoWithFields(msg interface{}, fields map[string]interface{}) {
	for _, log := range l.Logs {
		log.WithFields(fields).Info(msg)
	}
}

// InfoWithFields logs a message with level WARN and a map of fields to all loggers.
func (l *CompositeLogger) WarnWithFields(msg interface{}, fields map[string]interface{}) {
	for _, log := range l.Logs {
		log.WithFields(fields).Warn(msg)
	}
}
