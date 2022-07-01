# scanner_backend

## Description

scanner_backend is a backend and front end side [using go templates] of tg_scanner application

Application will always create a dir:
  - logs - for logs file

## Technology

Go, GoTemplates, PostgreSQL, Gorilla-Mux, Logrus, Go-Sqlmock, Testify, Golang-Migrate


## Installation

```bash
#Make sure that you have installed PostgreSQL on your machine

 git clone git@github.com:VladPetriv/scanner_backend.git

 cd scanner_backend 

 go mod download

```

## Before start

Please create .env file with this fields:
- POSTGRES_USER = PostgreSQL user
- POSTGRES_PASSWORD = PostgreSQL user password
- POSTGRES_HOST = PostgreSQL host
- POSTGRES_DB = PostgreSQL database name
- MIGRATIONS_PATH = Path to migrations for example:"file://./internal/store/migrations" 
- PORT = Bind address which server going to use

## Usage

Start an application locally:

```bash
 make run-l # Or you can use go run ./cmd/main.go
```

Start with docker-compose:

```bash
 make build # Build docker compoe

 make run # Up docker compose
```

Running test suite


```bash
 make mock

 make test
```

Watch demo version:

[Telegram Overflow](https://telegram-overflow.herokuapp.com/)


