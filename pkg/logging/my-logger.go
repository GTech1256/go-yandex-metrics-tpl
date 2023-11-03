package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
)

type MyLogger struct {
	Entry *logrus.Entry
}

var (
	// Получение пути до проекта
	rootpath, rootpatherr = os.Getwd()
)

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else if rootpatherr == nil {
		slash := strings.LastIndex(file, rootpath)
		if slash >= 0 {
			// Удаляет полный путь до проекта для вывода относительного пути
			file = file[len(rootpath):]
		}
	}
	// выводит относительный путь
	return fmt.Sprintf("%s:%d", file, line)
}

func NewMyLogger() *MyLogger {
	return &MyLogger{e}
}

func (m *MyLogger) GetLogger() *logrus.Entry {
	return m.Entry
}

func (m *MyLogger) WithFields(fields logrus.Fields) Logger {
	newEntry := m.Entry.WithFields(fields)
	newEntry.Data["file"] = fileInfo(2)
	return &MyLogger{Entry: newEntry}
}

func (m *MyLogger) WithField(key string, value interface{}) Logger {
	newEntry := m.Entry.WithField(key, value)
	newEntry.Data["file"] = fileInfo(2)
	return &MyLogger{Entry: newEntry}
}

func (m *MyLogger) Error(args ...interface{}) {
	m.Entry.Data["file"] = fileInfo(2)
	m.Entry.Error(args...)
}

func (m *MyLogger) Info(args ...interface{}) {
	m.Entry.Data["file"] = fileInfo(2)
	m.Entry.Info(args...)
}

func (m *MyLogger) Infof(format string, args ...interface{}) {
	m.Entry.Data["file"] = fileInfo(2)
	m.Entry.Infof(format, args...)
}

func (m *MyLogger) Errorf(format string, args ...interface{}) {
	m.Entry.Data["file"] = fileInfo(2)
	m.Entry.Errorf(format, args...)
}
