package sql

const (
	INSERT_COUNTER = "INSERT INTO counter(title, delta) VALUES($1, $2)"
	UPDATE_COUNTER = "UPDATE counter SET delta = $2 where title = $1"
	INSERT_GAUGE   = "INSERT INTO gauge(title, value) VALUES($1, $2)"
	UPDATE_GAUGE   = "UPDATE gauge SET value = $2 where title = $1"
)
