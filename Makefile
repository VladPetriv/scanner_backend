.PHONY:build

build:
	go build -o server ./cmd/main.go

.PHONY:run

run:
	go run ./cmd/main.go

