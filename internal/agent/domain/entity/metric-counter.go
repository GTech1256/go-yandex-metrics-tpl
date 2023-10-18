package entity

type CounterValue = int64

type MetricCounter struct {
	Type  Type
	Name  string
	Value CounterValue
}
