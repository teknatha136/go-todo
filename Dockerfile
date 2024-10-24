# Use an official Go image as the builder
FROM golang:1.22-alpine AS builder

# Install necessary build dependencies
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Cache dependencies by copying go.mod and go.sum separately
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app with optimizations and stripping debug information for a smaller binary
RUN go build -ldflags="-s -w" -o main .

# Use a minimal base image for the final container
FROM alpine:latest

# Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates

# Set the working directory inside the container
WORKDIR /app

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/main .

# Copy the static files from the builder stage
COPY --from=builder /app/static ./static

# Expose the port the app will run on (change if needed)
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
