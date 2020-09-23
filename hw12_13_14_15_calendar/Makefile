generate:
	go generate ./...

build:
	go build -o ./bin/calendar ./cmd/calendar
	go build -o ./bin/calendar_scheduler ./cmd/calendar_scheduler
	go build -o ./bin/calendar_sender ./cmd/calendar_sender

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || go install github.com/golangci/golangci-lint/cmd/golangci-lint

lint: install-lint-deps
	golangci-lint run ./...

up: build
	docker-compose -f deployments/docker-compose.yaml -p calendar up

down:
	docker-compose -f deployments/docker-compose.yaml -p calendar down

.PHONY: build test lint generate up down
