# Start from the Go image
FROM golang:1.23-alpine as builder

# Set the working directory
WORKDIR /app

# Copy Go modules and source files
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the application
RUN go build -o main ./cmd

# Start a minimal image for running the app
FROM debian:buster

# Set working directory and copy the binary
WORKDIR /app
COPY --from=builder /app/main .

# Copy the .env file into the container (make sure it's in the same directory as your Dockerfile)
COPY .env .env

# Expose the app port and run the app
EXPOSE 8080
CMD ["./main"]
