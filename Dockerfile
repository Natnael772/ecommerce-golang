# Step 1: Use the official Golang image as the build base
FROM golang:1.25 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Step 2: Copy the source code to the container
COPY . .

# Step 3: Build the Go app (targeting the main.go inside cmd/api/main)
WORKDIR /app/cmd/api

# Run the build command to compile the Go app
RUN go build -o /app/main .

# Step 4: Start a new stage to run the app from a clean image
FROM debian:bullseye-slim

# Install ca-certificates
RUN apt-get update && apt-get install -y ca-certificates

# Set the working directory inside the container (root directory)
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose the port your app is running on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
