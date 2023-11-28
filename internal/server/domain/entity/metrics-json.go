package entity

type MetricJSON struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (receiver MetricJSON) GetIsGauge() bool {
	return receiver.MType == "gauge"
}

func (receiver MetricJSON) GetIsCounter() bool {
	return receiver.MType == "counter"
}
