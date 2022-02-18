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

coverage.out: table/*.go go.*
	@go test -coverprofile=coverage.out ./table

.PHONY: fmt
fmt:
	@go fmt ./...

