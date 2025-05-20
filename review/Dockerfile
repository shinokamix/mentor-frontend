FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o review-service ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/review-service /app/review-service

COPY --from=builder /app/.env /app/.env

COPY --from=builder /app/keys /app/keys

EXPOSE 8082

CMD ["/app/review-service"]