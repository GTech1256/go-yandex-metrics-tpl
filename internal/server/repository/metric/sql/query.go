package sql

const (
	InsertCounter = "INSERT INTO counter(title, delta) VALUES($1, $2)"
	UpdateCounter = "UPDATE counter SET delta = $2 where title = $1"
	InsertGauge   = "INSERT INTO gauge(title, value) VALUES($1, $2)"
	UpdateGauge   = "UPDATE gauge SET value = $2 where title = $1"
)
