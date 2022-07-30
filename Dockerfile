FROM golang:1.18
ENV GOPATH=/

WORKDIR /src/
COPY ./ /src/
COPY ./.env ./src
COPY ./templates ./src/templates
COPY ./db/migrations ./src



#Build go application
RUN go mod download; CGO_ENABLED=0 go build -o /scanner-backend ./cmd/main.go


CMD ["./scanner-backend"]
