.PHONY: build test lint clean docker-build install

BINARY  := devops
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o bin/$(BINARY) ./

test:
	go test ./... -v -race -coverprofile=coverage.out

coverage: test
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/ coverage.out coverage.html

docker-build:
	docker build -t go-utils:$(VERSION) .

install:
	go install $(LDFLAGS) ./
