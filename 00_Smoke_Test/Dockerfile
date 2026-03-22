# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY main.go .

RUN go build -o server main.go

# Run stage (small image)
# This second FROM starts a clean slate, meaning whatever files we copied from host
# to our container in previous steps (ex: main.go in above stage) will not be carried over
# Even the base image is cleared
# Only those files that are copied again in following COPY instruction will exist in the final image
# Considering there are no more FROM instructions after this one.
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]