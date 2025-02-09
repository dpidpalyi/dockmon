FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/app/

FROM alpine

WORKDIR /app

COPY --from=builder /app/main .
COPY .app.env .

CMD ["/bin/sh", "-c", "source /app/.app.env && /app/main"]
