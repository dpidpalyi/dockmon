FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o pinger ./cmd/app/

FROM alpine

WORKDIR /app

COPY --from=builder /app/pinger .
COPY .env .

CMD ["./pinger"]
