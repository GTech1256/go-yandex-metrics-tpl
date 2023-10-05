package logging

import "github.com/sirupsen/logrus"

type Logger interface {
	WithFields(fields logrus.Fields) Logger
	WithField(key string, value interface{}) Logger
	Error(args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}
