include ./.env

.PHONY:build-l

build-l:
	go build -o server ./cmd/main.go

.PHONY:run-l

run-l:
	go run ./cmd/main.go

migrate_up:
	migrate -path ./db/migrations/ -database $(DATABASE_URL) -verbose up

.PHONY: migrate_down

migrate_down:
	migrate -path ./db/migrations/ -database $(DATABASE_URL) -verbose down

.PHONY: test

test:
	go test -v ./...

.PHONY: build

build:
	docker-compose build

.PHONY: run

run:
	docker-compose up

.PHONY: mock

mock:
	cd ./internal/store/; go generate;
	cd ./internal/service/; go generate;

