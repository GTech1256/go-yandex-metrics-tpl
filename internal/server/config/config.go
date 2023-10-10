package config

import (
	"flag"
	"os"
)

type Configurable interface {
	Load()
}

type Config struct {
	// port - Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
	Port *string
}

func NewConfig() Configurable {
	return &Config{}
}

func (c *Config) Load() {
	c.Port = flag.String("a", ":8080", "address and port to run server")
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		c.Port = &envRunAddr
	}

	flag.Parse()
}
