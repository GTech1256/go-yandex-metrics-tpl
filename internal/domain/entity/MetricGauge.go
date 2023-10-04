package entity

type GaugeValue = float64 // float64

type MetricGauge struct {
	Type  Type
	Name  string
	Value GaugeValue
}
