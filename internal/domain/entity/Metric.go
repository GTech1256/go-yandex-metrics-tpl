package entity

type Type string

// Возможные Типы метрики
const (
	NoType  Type = Type("")        // float64
	Gauge   Type = Type("gauge")   // float64
	Counter Type = Type("counter") // int64
)

// Возможные Значения метрики
type (
	GaugeValue   float64 // float64
	CounterValue int64   // int64
)

type Metric struct {
	Type  Type
	Name  string
	Value interface{}
}

type MetricGauge struct {
	Type  Type
	Name  string
	Value GaugeValue
}

type MetricCounter struct {
	Type  Type
	Name  string
	Value CounterValue
}

//type Name string
//type Value string
