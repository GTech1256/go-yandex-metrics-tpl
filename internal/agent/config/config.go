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
	c.ServerPort = flag.String("a", ":8080", "address and port to run server")
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		c.ServerPort = &envRunAddr
	}

	c.ReportInterval = flag.Int("r", 10, "frequency of sending metrics to the server")
	if envRunAddr := os.Getenv("REPORT_INTERVAL"); envRunAddr != "" {
		atoi, err := strconv.Atoi(envRunAddr)
		if err == nil {
			c.ReportInterval = &atoi
		}
	}

	c.PollInterval = flag.Int("p", 2, "the frequency of polling metrics")
	if envRunAddr := os.Getenv("POLL_INTERVAL"); envRunAddr != "" {
		atoi, err := strconv.Atoi(envRunAddr)
		if err == nil {
			c.PollInterval = &atoi
		}
	}

	flag.Parse()
}
