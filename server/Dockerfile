# Start from a minimal base image, such as Alpine Linux
FROM golang:alpine AS builder

# Set necessary environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory in the container
WORKDIR /root

# Copy the Go modules manifest and lock files
COPY . .

# Change directory to /app/api
WORKDIR /root/api

# Download dependencies
RUN go mod download

# Build the Go app
RUN go build -o app .

# Start a new stage from scratch
FROM alpine:latest

# Set necessary environment variables
ENV TZ=UTC

# Set the working directory in the container
WORKDIR /root/

# Copy the executable from the builder stage
COPY --from=builder /root/api/app .

# Copy GeoLite2-Country.mmdb
COPY --from=builder /root/GeoLite2-Country.mmdb .

# Expose port 3000 to the outside world
EXPOSE 3000

# Command to run the executable
CMD ["./app"]
