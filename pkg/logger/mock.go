package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

type LoggerMock struct {
	mock.Mock
}

func (m *LoggerMock) WithFields(fields logrus.Fields) Logger {
	m.Called(fields)
	return m
}

func (m *LoggerMock) WithField(key string, value interface{}) Logger {
	m.Called(key, value)
	return m
}

func (m *LoggerMock) Error(args ...interface{}) {
	fmt.Println(args...)
	m.Called(args...)
}

func (m *LoggerMock) Info(args ...interface{}) {
	fmt.Println(args...)
	m.Called(args...)
}

func (m *LoggerMock) Infof(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
	m.Called(append([]interface{}{format}, args...))
}

func (m *LoggerMock) Errorf(format string, args ...interface{}) {
	fmt.Println(fmt.Errorf(format, args...))
	m.Called(append([]interface{}{format}, args...))
}
