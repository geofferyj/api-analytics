# Start from a minimal base image, such as Alpine Linux
FROM golang:alpine AS builder

# Set necessary environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory in the container
WORKDIR /app

# Copy the Go modules manifest and lock files
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o app .

# Start a new stage from scratch
FROM alpine:latest

# Set necessary environment variables
ENV TZ=UTC

# Set the working directory in the container
WORKDIR /root/

# Copy the executable from the builder stage
COPY --from=builder /app/app .

# Expose port 8000 to the outside world
EXPOSE 8000

# Command to run the executable
CMD ["./app"]
