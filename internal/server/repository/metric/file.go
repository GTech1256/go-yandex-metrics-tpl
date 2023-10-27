package metric

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type MetricsJSON struct {
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
	file, err := os.OpenFile(fileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
		//logger.Error("Error:", err)
	}

	rw := bufio.NewReadWriter(bufio.NewReader(file), bufio.NewWriter(file))
	return &fileStorage{
		file: file,
		rw:   rw,
	}, nil
}

func (s *fileStorage) readLine() (*MetricsJSON, error) {
	bytes, err := s.rw.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	metric := &MetricsJSON{}
	err = json.Unmarshal(bytes, metric)
	if err != nil {
		return nil, err
	}

	return metric, nil
}

func (s *fileStorage) ReadAll() ([]*MetricsJSON, error) {
	metrics := make([]*MetricsJSON, 0)
	// Читаем и выводим весь файл
	for {
		line, err := s.readLine()
		if err != nil {
			break // Конец файла
		}
		fmt.Print(line)

		metrics = append(metrics, line)
	}

	return metrics, nil
}

func (s *fileStorage) Write(metric *MetricsJSON) error {
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
