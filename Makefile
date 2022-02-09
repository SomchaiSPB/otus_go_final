BIN := "./bin/previewer"

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/main.go

run:
	docker-compose up -d --build

stop:
	docker-compose down

restart:
	docker-compose restart

lint-fix:
	gofmt -s -w .
	golangci-lint run --fix

test:
	go test -race -count 10 ./test/...

lint: install-lint-deps
	golangci-lint run ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

.PHONY: all test clean build