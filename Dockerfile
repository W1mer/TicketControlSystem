FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o task-manager ./cmd/server

FROM alpine:latest
RUN apk add --no-cache libc6-compat
WORKDIR /app
COPY --from=builder /app/task-manager .
COPY .env .
CMD ["./task-manager"]