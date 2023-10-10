package config

import (
	"flag"
	"os"
)

type Configurable interface {
	Load()
}

type Config struct {
	// Port - Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
	Port *string
}

func NewConfig() Configurable {
	return &Config{}
}

func (c *Config) Load() {

	var (
		// Hack для тестирования
		command = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		port    = command.String("a", ":8080", "address and port to run server")
		portEnv = os.Getenv("ADDRESS")
	)

	c.Port = port
	if portEnv != "" {
		c.Port = &portEnv
	}

	// Тесты запускают несколько раз метод Load.
	// А несколько раз flag.Parse() нельзя вызывать
	// Из-за этого хак с командными флагами
	command.Parse(os.Args[1:])
}
