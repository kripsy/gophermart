# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.20

WORKDIR /app

COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
RUN go mod download

COPY ./bin ./bin
COPY ./cmd/auth ./cmd/auth
COPY ./db/auth ./db/auth
COPY ./docs/auth ./docs/auth
COPY ./internal/auth ./internal/auth
COPY ./*.md ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/auth ./cmd/auth/main.go
CMD ["./bin/auth"]

EXPOSE 8080