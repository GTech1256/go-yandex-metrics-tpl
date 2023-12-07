package config

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	// ServerPort - Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
	ServerPort *string

	// ReportInterval - Флаг -r=<ЗНАЧЕНИЕ> позволяет переопределять reportInterval — частоту отправки метрик на сервер (по умолчанию 10 секунд).
	ReportInterval *int

	// PollInterval - Флаг -p=<ЗНАЧЕНИЕ> позволяет переопределять pollInterval — частоту опроса метрик из пакета runtime (по умолчанию 2 секунды).
	PollInterval *int

	// Batch - Флаг -b=<ЗНАЧЕНИЕ> указывает отправлять метрику массивом JSON или отправлять каждую метрику отдельно (по умолчанию отдельно(false)).
	Batch *bool

	// HashKey - Флаг -k=<ЗНАЧЕНИЕ> При наличии ключа агент должен вычислять хеш и передавать в HTTP-заголовке запроса с именем HashSHA256.
	HashKey *string

	// RateLimit - Флаг -l=<ЗНАЧЕНИЕ> Ограничивает количество одновременно исходящих запросов на сервер «сверху».
	RateLimit *int
}

const EmptyStringFlagValue = ""
const EmptyIntFlagValue = 0

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load() {

	var (
		// Hack для тестирования
		command                                     = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		serverPort                                  = command.String("a", ":8080", "address and port to run memory")
		serverPortEnv, serverPortEnvPresent         = os.LookupEnv("ADDRESS")
		reportInterval                              = command.Int("r", 10, "frequency of sending metrics to the memory")
		reportIntervalEnv, reportIntervalEnvPresent = os.LookupEnv("REPORT_INTERVAL")
		pollInterval                                = command.Int("p", 2, "the frequency of polling metrics")
		pollIntervalEnv, pollIntervalEnvPresent     = os.LookupEnv("POLL_INTERVAL")
		batch                                       = command.Bool("b", false, "request by batch strategy")
		batchEnv, batchEnvPresent                   = os.LookupEnv("BATCH")
		hashKey                                     = command.String("k", EmptyStringFlagValue, "вычисляет хеш для подписи данных по ключу")
		hashKeyEnv, hashKeyEnvPresent               = os.LookupEnv("KEY")
		rateLimit                                   = command.Int("l", EmptyIntFlagValue, "лимит одновременно исходящих запросов на сервер")
		rateLimitEnv, rateLimitEnvPresent           = os.LookupEnv("RATE_LIMIT")
	)

	c.ServerPort = serverPort
	if serverPortEnvPresent {
		c.ServerPort = &serverPortEnv
	}

	c.ReportInterval = reportInterval
	if reportIntervalEnvPresent {
		atoi, err := strconv.Atoi(reportIntervalEnv)
		if err == nil {
			c.ReportInterval = &atoi
		}
	}

	c.PollInterval = pollInterval
	if pollIntervalEnvPresent {
		atoi, err := strconv.Atoi(pollIntervalEnv)
		if err == nil {
			c.PollInterval = &atoi
		}
	}

	c.Batch = batch
	if batchEnvPresent {
		batchBool, err := strconv.ParseBool(batchEnv)
		if err == nil {
			c.Batch = &batchBool
		}
	}

	if *hashKey != EmptyStringFlagValue {
		c.HashKey = hashKey
	}
	if hashKeyEnvPresent {
		c.HashKey = &hashKeyEnv
	}

	if *rateLimit != EmptyIntFlagValue {
		c.RateLimit = rateLimit
	}
	if rateLimitEnvPresent {
		atoi, err := strconv.Atoi(rateLimitEnv)
		if err == nil {
			c.RateLimit = &atoi
		}
	}

	// Тесты запускают несколько раз метод Load.
	// А несколько раз flag.Parse() нельзя вызывать
	// Из-за этого хак с командными флагами
	command.Parse(os.Args[1:])
}
