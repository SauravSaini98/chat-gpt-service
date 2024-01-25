# Use an official Go runtime as a parent image
FROM golang:1.21.6-alpine AS builder

WORKDIR /go/src/app

# Copy only the necessary files for dependency resolution
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go application
RUN GIN_MODE=release CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myapp .

# Stage 2: Create a minimal image with only the binary
FROM alpine:latest

WORKDIR /go/src/app

# Copy only the built binary from the previous stage
COPY --from=builder /go/src/app/myapp .

ENV GIN_MODE release
# Expose the port your application runs on
EXPOSE 8080

# Command to run your application
CMD ["./myapp"]
