# Stage 1: Build the Go application
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download all Go module dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o /app/main ./cmd/main.go

# Stage 2: Create a smaller image for the application
FROM alpine:latest

# Install ca-certificates, ffmpeg, yt-dlp, and Python
RUN apk --no-cache add ca-certificates ffmpeg python3 py3-pip yt-dlp

# Set the working directory inside the container
WORKDIR /root/

# Copy the Go binary from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Expose the port on which the service will run
EXPOSE 8081

# Command to run the application
CMD ["./main"]

