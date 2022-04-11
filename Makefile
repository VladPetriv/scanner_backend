.PHONY:build

build:
	go build -o server ./cmd/main.go

.PHONY:run

run:
	go run ./cmd/main.go

migrate_up:
	migrate -path ./internal/store/migrations/ -database "postgresql://vlad:admin@localhost:5432/tg_scanner?sslmode=disable" -verbose up

.PHONY: migrate_down

migrate_down:
	migrate -path ./internal/store/migrations/ -database "postgresql://vlad:admin@localhost:5432/tg_scanner?sslmode=disable" -verbose down

.PHONT: test

test:
	go test -v ./...
