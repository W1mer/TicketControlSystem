FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o task-manager ./cmd/server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/task-manager .
CMD ["./task-manager"]