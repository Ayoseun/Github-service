# Use the official Golang image as the base image
FROM golang:1.20-alpine as builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main cmd/main.go

# Use a minimal image for running the application
FROM alpine:3.18

# Set the working directory inside the container
WORKDIR /app

# Copy the built application from the builder image
COPY --from=builder /app/main .

# Expose the port on which the application will run
EXPOSE 8080

# Command to run the Go application
CMD ["./main"]
