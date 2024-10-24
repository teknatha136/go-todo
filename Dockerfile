FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Cache dependencies by copying go.mod and go.sum separately
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the Go app with optimizations and stripping debug information for a smaller binary
RUN go build -ldflags="-s -w" -o main .

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/static ./static

# Command to run the executable
CMD ["./main"]
