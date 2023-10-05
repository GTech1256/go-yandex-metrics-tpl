package logging

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

type LoggerMock struct {
	mock.Mock
}

func (m *LoggerMock) WithFields(fields logrus.Fields) *logrus.Entry {
	m.Called(fields)
	return &logrus.Entry{}
}

func (m *LoggerMock) WithField(key string, value interface{}) *logrus.Entry {
	m.Called(key, value)
	return &logrus.Entry{}
}

func (m *LoggerMock) Error(args ...interface{}) {
	m.Called(args...)
}

func (m *LoggerMock) Info(args ...interface{}) {
	m.Called(args...)
}

func (m *LoggerMock) Infof(format string, args ...interface{}) {
	m.Called(append([]interface{}{format}, args...))
}

func (m *LoggerMock) Errorf(format string, args ...interface{}) {
	m.Called(append([]interface{}{format}, args...))
}
