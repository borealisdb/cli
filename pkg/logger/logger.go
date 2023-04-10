package logger

import (
	"github.com/sirupsen/logrus"
)

var logLevels = map[string]logrus.Level{
	"info":    logrus.InfoLevel,
	"debug":   logrus.DebugLevel,
	"warning": logrus.WarnLevel,
}

func New(logLevelRaw string) *logrus.Entry {
	logLevel := logLevels[logLevelRaw]
	l := logrus.New()
	l.SetLevel(logLevel)
	logger := logrus.NewEntry(l)
	prefixedLogger := logger.WithField("component", "borealis-cli")
	return prefixedLogger
}
