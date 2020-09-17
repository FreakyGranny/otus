generate:
	go generate ./...

build:
	go build -o ./bin/calendar ./cmd/calendar

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || go install github.com/golangci/golangci-lint/cmd/golangci-lint

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build test lint generate
