FROM golang:latest AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o main ./cmd/api

# ------------------- Runtime -------------------
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY migrations ./migrations

# 🧠 Убедимся, что передаются аргументы из docker-compose
ENTRYPOINT ["/bin/sh", "-c", "./main $0 $@"]
CMD [""]
