package config

import (
	"flag"
	"os"
	"strconv"
)

type Configurable interface {
	Load()
}

type Config struct {
	// ServerPort - Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
	ServerPort *string

	// ReportInterval - Флаг -r=<ЗНАЧЕНИЕ> позволяет переопределять reportInterval — частоту отправки метрик на сервер (по умолчанию 10 секунд).
	ReportInterval *int

	// PollInterval - Флаг -p=<ЗНАЧЕНИЕ> позволяет переопределять pollInterval — частоту опроса метрик из пакета runtime (по умолчанию 2 секунды).
	PollInterval *int
}

func NewConfig() Configurable {
	return &Config{}
}

func (c *Config) Load() {

	var (
		// Hack для тестирования
		command           = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		serverPort        = command.String("a", ":8080", "address and port to run server")
		serverPortEnv     = os.Getenv("ADDRESS")
		reportInterval    = command.Int("r", 10, "frequency of sending metrics to the server")
		reportIntervalEnv = os.Getenv("REPORT_INTERVAL")
		pollInterval      = command.Int("p", 2, "the frequency of polling metrics")
		pollIntervalEnv   = os.Getenv("POLL_INTERVAL")
	)

	c.ServerPort = serverPort
	if serverPortEnv != "" {
		c.ServerPort = &serverPortEnv
	}

	c.ReportInterval = reportInterval
	if reportIntervalEnv != "" {
		atoi, err := strconv.Atoi(reportIntervalEnv)
		if err == nil {
			c.ReportInterval = &atoi
		}
	}

	c.PollInterval = pollInterval
	if pollIntervalEnv != "" {
		atoi, err := strconv.Atoi(pollIntervalEnv)
		if err == nil {
			c.PollInterval = &atoi
		}
	}

	// Тесты запускают несколько раз метод Load.
	// А несколько раз flag.Parse() нельзя вызывать
	// Из-за этого хак с командными флагами
	command.Parse(os.Args[1:])
}
