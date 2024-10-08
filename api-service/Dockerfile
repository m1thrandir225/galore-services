# Stage 1: Build the Go application
FROM golang:1.21-alpine AS builder

# Install necessary build tools
RUN apk update && apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules definition and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go binary
RUN go build -o main .

# Stage 2: Create a small image for the runtime
FROM alpine:latest

# Set up the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Copy the env file
COPY app.prod.env ./app.env

# Expose the port Gin will run on
EXPOSE 8080

# Command to run the Go binary
CMD ["./main"]

