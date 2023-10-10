package metricvalidator

import "errors"

var (
	// ErrNotCorrectURL возвращает ошибку, если формат url не подошел
	ErrNotCorrectURL   = errors.New("not correct url")
	ErrNotCorrectName  = errors.New("not correct metricName")
	ErrNotCorrectValue = errors.New("not correct metricValue")
	ErrNotCorrectType  = errors.New("not correct metricType")
)
