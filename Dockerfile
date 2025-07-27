FROM golang:1.23-alpine3.20 AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=1
RUN go build -o ./build/monotonic ./cmd/bot/main.go

# Lightweight final container
FROM alpine:latest

RUN apk add --no-cache sqlite-libs

WORKDIR /app

COPY --from=builder /app/build ./build

CMD ["./build/monotonic"]
