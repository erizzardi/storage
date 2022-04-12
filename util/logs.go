package util

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Wrapper around logrus.Logger and logrus.Entry.
// It has both logger.SetLayer() and entry.WithFields()
type Logger struct {
	logger *logrus.Logger
	fields logrus.Fields
}

// Wrapper arount logrus.New()
func NewLogger() *Logger {
	// return &Logger{logger: logrus.New()}

	ret := &logrus.Logger{
		Out: os.Stderr,
		Formatter: &logrus.TextFormatter{
			DisableColors:   false,
			ForceQuote:      false,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		},
	}

	return &Logger{logger: ret}
}

func (logger *Logger) GetFields() logrus.Fields { return logger.fields }

// Wrapper for logrus.SetLevel()
func (logger *Logger) SetLevel(level logrus.Level) {
	logger.logger.SetLevel(level)
}

func (logger *Logger) SetDefaultFields(fields logrus.Fields) {
	logger.fields = fields
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.logger.WithFields(logger.fields).Debug(args)
}

func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.logger.WithFields(logger.fields).Debugf(format, args)
}

func (logger *Logger) Info(args ...interface{}) {
	logger.logger.WithFields(logger.fields).Info(args)
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.logger.WithFields(logger.fields).Infof(format, args)
}

func (logger *Logger) Warn(args ...interface{}) {
	logger.logger.WithFields(logger.fields).Warn(args)
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.logger.WithFields(logger.fields).Warnf(format, args)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.logger.WithFields(logger.fields).Error(args)
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.logger.WithFields(logger.fields).Errorf(format, args)
}

func (logger *Logger) Fatal(args ...interface{}) {
	logger.logger.WithFields(logger.fields).Fatal(args)
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	logger.logger.WithFields(logger.fields).Fatalf(format, args)
}

func InitLogger(logger *Logger, level string, fields logrus.Fields) {
	switch level {
	case "DEBUG":
		logger.logger.Level = logrus.DebugLevel
	case "INFO":
		logger.logger.Level = logrus.InfoLevel
	case "WARN":
		logger.logger.Level = logrus.WarnLevel
	case "ERROR":
		logger.logger.Level = logrus.ErrorLevel
	case "FATAL":
		logger.logger.Level = logrus.FatalLevel
	}

	logger.SetDefaultFields(fields)
}
