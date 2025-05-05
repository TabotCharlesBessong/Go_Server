# .PHONY: build run test clean migrate

# Build the application
build:
	go build -o bin/server cmd/main.go

# Run the application
run:
	go run cmd/main.go

# Run migrations
# migrate:
# 	go run cmd/migrate/main.go

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out  $@, $(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Install dependencies
deps:
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Default target
all: deps fmt vet test build