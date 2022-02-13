.PHONY: default
default:
	@go run examples/dimensions/main.go

.PHONY: test
test:
	@go test ./table

.PHONY: fmt
fmt:
	@go fmt ./...

