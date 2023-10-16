package logging

import "github.com/sirupsen/logrus"

type MyLogger struct {
	entry *logrus.Entry
}

func NewMyLogger() *MyLogger {
	return &MyLogger{e}
}

func (m *MyLogger) WithFields(fields logrus.Fields) Logger {
	newEntry := m.entry.WithFields(fields)
	return &MyLogger{entry: newEntry}
}

func (m *MyLogger) WithField(key string, value interface{}) Logger {
	newEntry := m.entry.WithField(key, value)
	return &MyLogger{entry: newEntry}
}

func (m *MyLogger) Error(args ...interface{}) {
	m.entry.Error(args...)
}

func (m *MyLogger) Info(args ...interface{}) {
	m.entry.Info(args...)
}

func (m *MyLogger) Infof(format string, args ...interface{}) {
	m.entry.Infof(format, args...)
}

func (m *MyLogger) Errorf(format string, args ...interface{}) {
	m.entry.Errorf(format, args...)
}
