FROM golang:1.17.1 AS build

ENV GOPATH=/

WORKDIR /src/
COPY ./ /src/

#Build go application
RUN go mod download; CGO_ENABLED=0 go build -o /scanner-backend ./cmd/main.go

#Install postgresql
FROM alpine:latest
COPY --from=build /scanner-backend /scanner-backend

COPY ./.env ./
COPY ./templates ./templates
COPY ./wait-for-postgres.sh ./
COPY ./internal/store/migrations ./

RUN apk --no-cache add postgresql-client && chmod +x wait-for-postgres.sh

CMD ["./scanner-backend"]
