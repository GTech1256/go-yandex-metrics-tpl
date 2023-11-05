package sql

import (
	"context"
	entity2 "github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/entity"
	"github.com/GTech1256/go-yandex-metrics-tpl/internal/server/domain/metric"
	"github.com/jackc/pgx/v5"
)

//var (
//	gauge   = make(map[string]float64)
//	counter = make(map[string]int64)
//)

//type MemStorage struct {
//	gauge   map[string]float64
//	counter map[string]int64
//}

type storage struct {
	db *pgx.Conn
	//memStorage MemStorage
}

func NewStorage(db *pgx.Conn) *storage {
	//memStorage := MemStorage{
	//	gauge:   gauge,
	//	counter: counter,
	//}

	return &storage{
		db: db,
		//memStorage: memStorage,
	}
}

// SaveGauge новое значение должно замещать предыдущее.
func (s *storage) SaveGauge(ctx context.Context, gauge *entity2.MetricGauge) error {
	//s.memStorage.gauge[gauge.Name] = float64(gauge.Value)
	//
	//fmt.Printf("%v %+v \n", len(s.memStorage.gauge), s.memStorage)
	//
	//return nil

	return nil
}

// SaveCounter новое значение должно добавляться к предыдущему, если какое-то значение уже было известно серверу.
func (s *storage) SaveCounter(ctx context.Context, counter *entity2.MetricCounter) error {
	//if _, isOk := s.memStorage.counter[counter.Name]; !isOk {
	//	s.memStorage.counter[counter.Name] = int64(counter.Value)
	//} else {
	//	s.memStorage.counter[counter.Name] += int64(counter.Value)
	//}
	//
	//fmt.Printf("%v %+v \n", len(s.memStorage.counter), s.memStorage)
	//
	//return nil

	return nil
}

// GetGauge - возвращает значение Gauge из хранилища
func (s *storage) GetGaugeValue(name string) (*entity2.GaugeValue, error) {
	//value, ok := s.memStorage.gauge[name]
	//if ok {
	//	return &value, nil
	//}
	//
	//return nil, nil

	return nil, nil
}

// GetCounter - возвращает значение Counter из хранилища
func (s *storage) GetCounterValue(name string) (*entity2.CounterValue, error) {
	//value, ok := s.memStorage.counter[name]
	//if ok {
	//	return &value, nil
	//}
	//
	//return nil, nil

	return nil, nil
}

func (s *storage) GetAllMetrics() *metric.AllMetrics {
	//return &metric.AllMetrics{
	//	Gauge:   s.memStorage.gauge,
	//	Counter: s.memStorage.counter,
	//}

	return nil
}

func (s *storage) Ping(ctx context.Context) error {
	return s.db.Ping(ctx)
}
