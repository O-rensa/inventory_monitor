OUTPUT=inventory_monitor

build:
	go build -o ./bin/${OUTPUT} ./cmd/main.go

run:
	go build -o ./bin/${OUTPUT} ./cmd/main.go
	./bin/${OUTPUT}

test:
	@go test -v ./...

tidy:
	go mod tidy

migrate:
	go run ./pkg/db/main.go
