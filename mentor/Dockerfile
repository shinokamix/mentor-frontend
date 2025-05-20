FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o mentor-service ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/mentor-service /app/mentor-service

COPY --from=builder /app/.env /app/.env

EXPOSE 8083 50051

CMD ["/app/mentor-service"]