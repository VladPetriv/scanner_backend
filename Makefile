.PHONY:build

build:
	go build -o server ./cmd/main.go

.PHONY:run

run:
	go run ./cmd/main.go

migrate_up:
	migrate -path ./internal/store/migrations/ -database "postgresql://vlad:admin@localhost:5432/scanner?sslmode=disable" -verbose up

.PHONY: migrate_down

migrate_down:
	migrate -path ./internal/store/migrations/ -database "postgresql://vlad:admin@localhost:5432/scanner?sslmode=disable" -verbose down

.PHONY: test

test:
	go test -v ./...

.PHONY: docker_build

docker_build:
	docker-compose build

.PHONY: docker_run

docker_run:
	docker-compose up


