package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	// Port - Флаг -a=<ЗНАЧЕНИЕ> отвечает за адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
	Port *string

	StoreInterval time.Duration

	FileStoragePath *string

	Restore *bool
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load() {

	var (
		// Hack для тестирования
		command            = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		port               = command.String("a", ":8080", "address and port to run metric")
		portEnv            = os.Getenv("ADDRESS")
		storeInterval      = command.Int("i", 300, "the number of seconds after which the metric is saved to disk")
		storeIntervalEnv   = os.Getenv("STORE_INTERVAL")
		fileStoragePath    = command.String("f", "/tmp/metrics-db.json", "the path to the file for saving metrics")
		fileStoragePathEnv = os.Getenv("FILE_STORAGE_PATH")
		restore            = command.Bool("r", true, "the path to the file for saving metrics")
		restoreEnv         = os.Getenv("RESTORE")
	)

	fmt.Println("storeIntervalEnv:", storeIntervalEnv)
	fmt.Println(os.Environ())

	c.Port = port
	if portEnv != "" {
		c.Port = &portEnv
	}

	c.StoreInterval = time.Duration(*storeInterval) * time.Second
	if storeIntervalEnv != "" {
		atoi, err := strconv.Atoi(storeIntervalEnv)
		if err == nil {
			c.StoreInterval = time.Duration(atoi) * time.Second
		}
	}

	c.FileStoragePath = fileStoragePath
	if fileStoragePathEnv != "" {
		c.FileStoragePath = &fileStoragePathEnv
	}

	c.Restore = restore
	if fileStoragePathEnv != "" {
		restoreEnvBool, err := strconv.ParseBool(restoreEnv)
		if err == nil {
			c.Restore = &restoreEnvBool
		}
	}

	// Тесты запускают несколько раз метод Load.
	// А несколько раз flag.Parse() нельзя вызывать
	// Из-за этого хак с командными флагами
	command.Parse(os.Args[1:])
}
