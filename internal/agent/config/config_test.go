package config

import (
	"os"
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Резервное копирование существующих аргументов командной строки и переменных окружения.
	oldArgs := os.Args
	oldEnv := os.Environ()
	defer func() {
		os.Args = oldArgs
		os.Clearenv()
		for _, e := range oldEnv {
			pair := strings.SplitN(e, "=", 2)
			os.Setenv(pair[0], pair[1])
		}
	}()

	// Инициализация аргументов командной строки для разбора флагов.
	os.Args = []string{"test", "-a=:8081", "-r=20", "-p=3"}

	// Установка переменных окружения для тестирования.
	os.Setenv("ADDRESS", ":8082")
	os.Setenv("REPORT_INTERVAL", "15")
	os.Setenv("POLL_INTERVAL", "5")

	// Создание объекта Config и загрузка его значений.
	config := NewConfig().(*Config)
	config.Load()

	if *config.ServerPort != ":8082" {
		t.Errorf("Ожидалось, что ServerPort будет ':8082', но получено: %s", *config.ServerPort)
	}

	if *config.ReportInterval != 15 {
		t.Errorf("Ожидалось, что ReportInterval будет 15, но получено: %d", *config.ReportInterval)
	}

	if *config.PollInterval != 5 {
		t.Errorf("Ожидалось, что PollInterval будет 5, но получено: %d", *config.PollInterval)
	}
}

func TestLoadConfigWithNoEnvVars(t *testing.T) {
	//command := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	// Резервное копирование существующих аргументов командной строки и переменных окружения.
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	// Инициализация аргументов командной строки для разбора флагов.
	os.Args = []string{"test", "-a=:8081", "-r=20", "-p=3"}

	// Создание объекта Config и загрузка его значений.
	config := NewConfig().(*Config)
	config.Load()

	if *config.ServerPort != ":8081" {
		t.Errorf("Ожидалось, что ServerPort будет ':8081', но получено: %s", *config.ServerPort)
	}

	if *config.ReportInterval != 20 {
		t.Errorf("Ожидалось, что ReportInterval будет 20, но получено: %d", *config.ReportInterval)
	}

	if *config.PollInterval != 3 {
		t.Errorf("Ожидалось, что PollInterval будет 3, но получено: %d", *config.PollInterval)
	}
}
