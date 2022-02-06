build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/main.go

run:
	docker-compose up -d --build

stop:
	docker-compose down

restart:
	docker-compose restart

test:
	go test -race -count 10 ./test/...

lint:
	golangci-lint run ./...