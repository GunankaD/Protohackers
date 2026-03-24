# Builder stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN go mod init budgetchat
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]