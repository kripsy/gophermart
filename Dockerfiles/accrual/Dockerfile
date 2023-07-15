# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.20

WORKDIR /app

COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
RUN go mod download

COPY ./bin ./bin
COPY ./cmd/accrual ./cmd/accrual
COPY ./db/accrual ./db/accrual
COPY ./docs/accrual ./docs/accrual
COPY ./internal/accrual ./internal/accrual
COPY ./*.md ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/accrual ./cmd/accrual/main.go
CMD ["./bin/accrual"]

EXPOSE 8080