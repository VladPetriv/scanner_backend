# scanner_backend

scanner_backend is a backend side represented as full-stack application with go-template for processing data from tg scanner


## Tech Stack

**Server:** 
- gorilla/mux
- Apache Kafka
- Go Templates

**DB:**
- PostgreSQL
- golang-migrate/migrate

**Testing:**
- stretchr/testify
- DATA-DOG/go-sqlmock
- mockery/mockery



## Environment Variables

To run this project, you will need to add the following environment variables to your ".env" file:
- `POSTGRES_USER` - PostgreSQL user
- `POSTGRES_PASSWORD` - PostgreSQL user password
- `POSTGRES_HOST` - PostgreSQL host
- `POSTGRES_DB` - PostgreSQL database name
- `MIGRATIONS_PATH` - Path to migrations:“file://./db/migrations”
- `PORT` - Bind address which server will use
- `DATABASE_URL` - this field you can use if you don’t want to create PostgreSQL fields

## Run Locally

Clone the project

```bash
   git clone git@github.com:VladPetriv/scanner_backend.git
```

Go to the project directory

```bash
  cd scanner_backend
```

Install dependencies

```bash
  go mod download
```

Start the server locally:

```bash
  # Make sure that Apache Kafka and PostgreSQL are running
  make run-l # Or you can use "go run ./cmd/main.go"
```

Start the server with docker compose:

```bash
  make build
  
  make run 
```
## Running Tests

To run tests, run the following command:

```bash
  # Run it only if "mocks" folder not exist or if you updated "service.go" or "repository.go" files
  make mock 
```

```bash
  make test 
```
## License

[MIT](https://choosealicense.com/licenses/mit/)
