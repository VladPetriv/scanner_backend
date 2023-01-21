include ./.env

.PHONY:build
build:
	go build -o server ./cmd/main.go

.PHONY:run
run:
	go run ./cmd/main.go

.PHONY: migrate_up
migrate_up:
	migrate -path ./db/migrations/ -database $(DATABASE_URL) -verbose up

.PHONY: migrate_down
migrate_down:
	migrate -path ./db/migrations/ -database $(DATABASE_URL) -verbose down

.PHONY: test
test:
	go test -v ./...

.PHONY: mock
mock:
	cd ./internal/store/; go generate;
	cd ./internal/service/; go generate;

.PHONY: lint
lint:
	golangci-lint run ./...
