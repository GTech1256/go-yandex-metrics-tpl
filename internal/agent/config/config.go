package config

import "flag"

type Configurable interface {
	Load()
}

type Config struct {
	// port - Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
	Port *string

	// reportInterval - Флаг -r=<ЗНАЧЕНИЕ> позволяет переопределять reportInterval — частоту отправки метрик на сервер (по умолчанию 10 секунд).
	ReportInterval *int

	// pollInterval - Флаг -p=<ЗНАЧЕНИЕ> позволяет переопределять pollInterval — частоту опроса метрик из пакета runtime (по умолчанию 2 секунды).
	PollInterval *int
}

func NewConfig() Configurable {
	return &Config{}
}

func (c *Config) Load() {
	c.Port = flag.String("a", ":8080", "address and port to run server")
	c.ReportInterval = flag.Int("r", 10, "frequency of sending metrics to the server")
	c.PollInterval = flag.Int("p", 2, "the frequency of polling metrics")

	flag.Parse()
}
