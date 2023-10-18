package entity

type Type string

// Возможные Типы метрики
const (
	NoType  Type = Type("")
	Gauge   Type = Type("gauge")   // float64
	Counter Type = Type("counter") // int64
)
