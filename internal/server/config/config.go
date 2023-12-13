package config

import (
	"flag"
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

	DatabaseDSN *string

	// HashKey - Флаг -k=<ЗНАЧЕНИЕ> При наличии ключа агент должен вычислять хеш и передавать в HTTP-заголовке запроса с именем HashSHA256.
	HashKey *string
}

const EmptyHashKey = ""

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load() {
	var (
		// Hack для тестирования
		command                                       = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		port                                          = command.String("a", ":8080", "address and port to run memory")
		portEnv, portEnvPresent                       = os.LookupEnv("ADDRESS")
		storeInterval                                 = command.Int("i", 300, "the number of seconds after which the memory is saved to disk")
		storeIntervalEnv, storeIntervalEnvPresent     = os.LookupEnv("STORE_INTERVAL")
		fileStoragePath                               = command.String("f", "/tmp/metrics-db.json", "the path to the file for saving metrics")
		fileStoragePathEnv, fileStoragePathEnvPresent = os.LookupEnv("FILE_STORAGE_PATH")
		restore                                       = command.Bool("r", true, "the path to the file for saving metrics")
		restoreEnv, restoreEnvPresent                 = os.LookupEnv("RESTORE")
		databaseDSN                                   = command.String("d", "", "the path to database connection")
		databaseDSNEnv, databaseDSNEnvPresent         = os.LookupEnv("DATABASE_DSN")
		hashKey                                       = command.String("k", EmptyHashKey, "вычисляет хеш для подписи данных по ключу")
		hashKeyEnv, hashKeyEnvPresent                 = os.LookupEnv("KEY")
	)

	c.Port = port
	if portEnvPresent {
		c.Port = &portEnv
	}

	c.StoreInterval = time.Duration(*storeInterval) * time.Second
	if storeIntervalEnvPresent {
		atoi, err := strconv.Atoi(storeIntervalEnv)
		if err == nil {
			c.StoreInterval = time.Duration(atoi) * time.Second
		}
	}

	c.FileStoragePath = fileStoragePath
	if fileStoragePathEnvPresent {
		c.FileStoragePath = &fileStoragePathEnv
	}

	c.Restore = restore
	if restoreEnvPresent {
		restoreEnvBool, err := strconv.ParseBool(restoreEnv)
		if err == nil {
			c.Restore = &restoreEnvBool
		}
	}

	c.DatabaseDSN = databaseDSN
	if databaseDSNEnvPresent {
		c.DatabaseDSN = &databaseDSNEnv
	}
	
	if *hashKey != EmptyHashKey {
		c.HashKey = hashKey
	}
	if hashKeyEnvPresent {
		c.HashKey = &hashKeyEnv
	}

	// Тесты запускают несколько раз метод Load.
	// А несколько раз flag.Parse() нельзя вызывать
	// Из-за этого хак с командными флагами
	command.Parse(os.Args[1:])
}

// пустое значение отключает функцию записи на диск
func (c Config) GetIsEnabledFileWrite() bool {
	return *c.FileStoragePath != ""
}

// пустое значение отключает функцию записи в SQL DB
func (c Config) GetIsEnabledSQLStore() bool {
	return *c.DatabaseDSN != ""
}
