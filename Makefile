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
	@go test ./table

.PHONY: fmt
fmt:
	@go fmt ./...

