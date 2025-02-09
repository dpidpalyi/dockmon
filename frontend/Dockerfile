FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o frontend ./cmd/web

FROM alpine

WORKDIR /app

COPY --from=builder /app/frontend .
COPY .env .
COPY ui ui

CMD ["./frontend"]
