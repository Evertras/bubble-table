.PHONY: default
default:
	go run examples/features/main.go

.PHONY: test
test:
	go test ./table

fmt:
	@go fmt ./...
