FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o api ./cmd/api

FROM alpine

WORKDIR /app

COPY --from=builder /app/api .
COPY .env .
COPY migrations ./migrations

CMD ["./api"]
