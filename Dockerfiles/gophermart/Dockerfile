# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.20

WORKDIR /app

COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
RUN go mod download

RUN mkdir -p -m 0600 ./bin
COPY ./cmd/gophermart ./cmd/gophermart
COPY ./db/gophermart ./db/gophermart
COPY ./docs/gophermart ./docs/gophermart
COPY ./internal/gophermart ./internal/gophermart
COPY ./internal/common ./internal/common
COPY ./*.md ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/gophermart ./cmd/gophermart/main.go
CMD ["./bin/gophermart"]

EXPOSE 8080