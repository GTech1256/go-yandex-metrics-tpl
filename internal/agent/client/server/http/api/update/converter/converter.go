package converter

import (
	"errors"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/client/server/dto"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/agent/domain/entity"
	"strconv"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

var (
	ErrNotValidType = errors.New("not valid type")
)

func UpdateDTOToMetrics(update *dto.Update) (*Metrics, error) {
	m := &Metrics{
		ID:    update.Name,
		MType: update.Type,
		//Delta: nil,
		//Value: nil,
	}

	switch update.Type {
	case string(entity.Gauge):
		floatValue, err := strconv.ParseFloat(update.Value, 64)
		if err != nil {
			return nil, err
		}
		m.Value = &floatValue
		return m, nil

	case string(entity.Counter):
		intValue, err := strconv.ParseInt(update.Value, 10, 64)
		if err != nil {
			return nil, err
		}
		m.Delta = &intValue
		return m, nil
	}

	return nil, ErrNotValidType
}
