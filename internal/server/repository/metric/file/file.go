package file

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/GTech1256/go-yandex-metrics-tpl/pkg/retry"
	"os"
)

type MetricJSON struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type fileStorage struct {
	file *os.File
	rw   *bufio.ReadWriter
}

func NewFileStorage(fileStoragePath string) (*fileStorage, error) {
	file, err := openFileRetry(fileStoragePath)
	if err != nil {
		return nil, err
	}

	rw := bufio.NewReadWriter(bufio.NewReader(file), bufio.NewWriter(file))
	return &fileStorage{
		file: file,
		rw:   rw,
	}, nil
}

func openFileRetry(fileStoragePath string) (*os.File, error) {
	var file *os.File
	var err error

	err = retry.MakeRetry(
		func() error {
			file, err = os.OpenFile(fileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

			// Когда файл заблокирован другим процессом
			// делается еще попытка
			if errors.Is(err, os.ErrExist) {
				return err
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return file, err
}

func (s *fileStorage) readLine() (*MetricJSON, error) {
	bytes, err := s.rw.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	metric := &MetricJSON{}
	err = json.Unmarshal(bytes, metric)
	if err != nil {
		return nil, err
	}

	return metric, nil
}

func (s *fileStorage) ReadAll() ([]*MetricJSON, error) {
	metrics := make([]*MetricJSON, 0)
	// Читаем и выводим весь файл
	for {
		line, err := s.readLine()
		if err != nil {
			break // Конец файла
		}

		metrics = append(metrics, line)
	}

	return metrics, nil
}

func (s *fileStorage) Write(metric *MetricJSON) error {
	bytes, err := json.Marshal(metric)
	if err != nil {
		return err
	}

	bytes = append(bytes, '\n')

	_, err = s.rw.Write(bytes)
	if err != nil {
		return err
	}

	err = s.rw.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (s *fileStorage) Truncate() error {
	// Укорачиваем файл до нулевой длины, что очистит его содержимое
	err := s.file.Truncate(0)
	if err != nil {
		return err
	}

	// Очищаем буфер и записываем изменения на диск
	return s.rw.Flush()
}
