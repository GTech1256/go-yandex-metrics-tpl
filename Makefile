build-server:
	go build -o cmd/server -v cmd/server/main.go
	mv cmd/server/main cmd/server/server

build-agent:
	go build -o cmd/agent -v cmd/agent/main.go
	mv cmd/agent/main cmd/agent/agent

.PHONY: run
build: build-server build-agent


start-agent:
	cmd/agent/agent

start-server:
	cmd/server/server

.PHONY: run
run:
	air

test:
	go test -count=1 ./...

test-coverage:
	go test -coverprofile cover.out ./...

test-coverage-html: test-coverage
	go tool cover -html=cover.out

test-autotests-server: build-server
	autotests/metricstest-darwin-arm64 -test.v -test.run=^TestIteration1$$ -binary-path=cmd/server/server

test-autotests-agent: build-server
	autotests/metricstest-darwin-arm64 -test.v -test.run=^TestIteration1$$ -binary-path=cmd/agent/agent



.DEFAULT_GOAL := run
