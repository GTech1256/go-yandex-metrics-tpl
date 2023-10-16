package metric

type AllMetrics struct {
	Gauge   map[string]float64
	Counter map[string]int64
}
