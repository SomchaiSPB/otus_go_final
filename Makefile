build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/main.go

run:
	docker-compose up -d --build

test:
	go test -race -count 10 ./test/...

lint:
	golangci-lint run ./...