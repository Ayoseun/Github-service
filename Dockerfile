# Use the official Golang image as the base image
FROM golang:1.20-alpine AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go Modules files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Use a minimal image to run the compiled Go app
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the Go app from the builder container
COPY --from=builder /app/main .

# Expose the port that the app runs on
EXPOSE 8080

# Command to run the Go app
CMD ["./main"]
