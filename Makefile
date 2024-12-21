.PHONY: run
run: build
	@./bin/main --config=config.yaml

.PHONY: build
build:
	@go build -o ./bin/main ./cmd/homeset

.PHONY: generate
generate:
	@go generate tools.go
	@go generate ./...

.PHONY: test
test:
	@go test -count=1 -p=4 -timeout=30s -race -cover -covermode=atomic -coverprofile=cover.out ./internal/...

	@go tool cover -html cover.out -o cover.html
