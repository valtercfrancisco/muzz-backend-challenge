# Start with the official Golang image as the base image for building
FROM golang:1.22.5-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code from the current directory to the Working Directory inside the container
COPY . .

# Build the Go application and output the executable to `/app/muzz-backend-challenge`
RUN go build -o /app/muzz-backend-challenge ./cmd/server

# Start a new stage from Alpine for a minimal runtime environment
FROM alpine:latest

# Install CA certificates to support HTTPS connections
RUN apk --no-cache add ca-certificates

# Set the working directory inside the container
WORKDIR /root/

# Copy the compiled binary from the previous stage into the final image
COPY --from=builder /app/muzz-backend-challenge .
COPY --from=builder /app/internal/db/migrations /db/migrations
COPY --from=builder /app/internal/db/mock /app/internal/db/mock
COPY db-variables.env .

# Ensure the binary is executable
RUN chmod +x ./muzz-backend-challenge

# Expose port 50051 to the outside world
EXPOSE 50051

# Command to run the executable
CMD ["./muzz-backend-challenge"]
