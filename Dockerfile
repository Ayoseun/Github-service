# Use the official Golang image as the base image for building
FROM golang:1.23-alpine as builder

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

# Create and set up the .env file
RUN echo "BASE_URL=https://api.github.com/repos" > .env \
    && echo "GITHUB_TOKEN=" >> .env \
    && echo "DEFAULT_OWNER=chromium" >> .env \
    && echo "DEFAULT_REPO=chromium" >> .env \
    && echo "BEGIN_FETCH_DATE=2023-01-01T00:00:00Z" >> .env \
    && echo "PORT=8080" >> .env \
    && echo "POSTGRES_USER=postgres" >> .env \
    && echo "POSTGRES_PASSWORD=Yourp@sswoird" >> .env \
    && echo "POSTGRES_HOST=db" >> .env \
    && echo "POSTGRES_DB=github_test" >> .env \
    && echo "POLL_INTERVAL=2" >> .env \
    && echo "PER_PAGE=100" >> .env

# Expose the port on which the application will run
EXPOSE 8080

# Command to run the Go application
CMD ["./main"]