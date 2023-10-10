package config

import "flag"

type Configurable interface {
	Load()
}

type Config struct {
	// AgentPort - Флаг -p=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-agent (по умолчанию localhost:8081).
	AgentPort *string

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
	c.AgentPort = flag.String("port", ":8081", "address and port to run agent")
	c.ServerPort = flag.String("a", ":8080", "address and port to run server")
	c.ReportInterval = flag.Int("r", 10, "frequency of sending metrics to the server")
	c.PollInterval = flag.Int("p", 2, "the frequency of polling metrics")

	flag.Parse()
}
