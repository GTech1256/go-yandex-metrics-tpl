build-server:
	go build -o cmd/server -v cmd/server/main.go
	mv cmd/server/main cmd/server/server

build-agent:
	go build -o cmd/agent -v cmd/agent/main.go
	mv cmd/agent/main cmd/agent/agent

build: build-server build-agent


start-agent:
	cmd/agent/agent

start-server:
	cmd/server/server

.PHONY: run
run:
	air

run-server:
	air

run-agent:
	go run cmd/agent/main.go

test:
	go test -count=1 ./...

test-coverage:
	go test -coverprofile cover.out ./...

test-coverage-html: test-coverage
	go tool cover -html=cover.out

test-coverage-p: test-coverage
	go tool cover -func cover.out | fgrep total | awk '{print $3}'

test-autotests-server: build-server
	autotests/metricstest-darwin-arm64 -test.v -test.run=^TestIteration1$$ -binary-path=cmd/server/server

test-autotests-agent: build-server
	autotests/metricstest-darwin-arm64 -test.v -test.run=^TestIteration1$$ -binary-path=cmd/agent/agent

db-migration-create:
	migrate create -ext sql -dir db/migrations -seq create_users_table

db-migration-up:
	export POSTGRESQL_URL="postgres://postgres:postgres@localhost:5432/yandex_metrics?sslmode=disable"
	migrate -database ${POSTGRESQL_URL} -path internal/server/config/db/migrations up

db-migration-down:
	export POSTGRESQL_URL="postgres://postgres:postgres@localhost:5432/yandex_metrics?sslmode=disable"
	migrate -database ${POSTGRESQL_URL} -path internal/server/config/db/migrations down

.DEFAULT_GOAL := run
