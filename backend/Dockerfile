# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o main .

# Final stage
FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/main .

# Copy keys
COPY --from=builder /app/keys ./keys

EXPOSE 8080 9999

CMD ["./main"]
