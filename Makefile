.PHONY:build-l

build-l:
	go build -o server ./cmd/main.go

.PHONY:run-l

run-l:
	go run ./cmd/main.go

migrate_up:
	migrate -path ./internal/store/migrations/ -database "postgresql://vlad:admin@localhost:5432/scanner?sslmode=disable" -verbose up

.PHONY: migrate_down

migrate_down:
	migrate -path ./internal/store/migrations/ -database "postgresql://vlad:admin@localhost:5432/scanner?sslmode=disable" -verbose down

.PHONY: test

test:
	go test -v ./...

.PHONY: build

build:
	docker-compose build

.PHONY: run

run:
	docker-compose up


