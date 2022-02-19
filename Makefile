.PHONY: features
features:
	@go run ./examples/features/main.go

.PHONY: dimensions
dimensions:
	@go run ./examples/dimensions/main.go

.PHONY: updates
updates:
	@go run ./examples/updates/*.go

.PHONY: test
test:
	@go test -race -cover ./table

.PHONY: test-coverage
test-coverage: coverage.out
	@go tool cover -html=coverage.out

.PHONY: lint
lint: ./bin/golangci-lint
	@./bin/golangci-lint run ./table

coverage.out: table/*.go go.*
	@go test -coverprofile=coverage.out ./table

.PHONY: fmt
fmt:
	@go fmt ./...

./bin/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.44.2

