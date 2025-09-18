# Stage 1
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY server/go.mod server/go.sum ./
RUN go mod download

COPY server/. .
RUN go build -o cms-server ./cmd/api

# Stage 2
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/cms-server .

EXPOSE 8080
CMD ["./cms-server"]
