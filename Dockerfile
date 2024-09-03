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
RUN echo "DATABASE_DEV_URL=postgres://postgres:Yourp@sswoird@db:5432/github_test" > .env \
    && echo "DATABASE_PROD_URL=postgresql://github_test_user:MfZyIf6vkrVm0O4YylZi3ig3MG3TV4VF@dpg-cr2kribtq21c73fa9kk0-a.oregon-postgres.render.com/github_test" >> .env \
    && echo "BASE_URL=https://api.github.com/repos" >> .env \
    && echo "SEED_REPO_OWNER=chromium" >> .env \
    && echo "SEED_REPO_NAME=chromium" >> .env \
    && echo "BEGIN_FETCH_DATE=2022-12-09T00:00:00Z" >> .env

# Expose the port on which the application will run
EXPOSE 8080

# Command to run the Go application
CMD ["./main"]
