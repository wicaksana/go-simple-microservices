# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
# Ensure you have go.mod and go.sum files in your project directory
# Run 'go mod init yourprojectname' and 'go mod tidy' if you don't
COPY go.mod ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download && go mod verify

# Copy the source code into the container
COPY . .

# Copy the templates directory
COPY templates ./templates

# Build the Go app
# CGO_ENABLED=0 for a statically linked binary
# GOOS=linux for cross-compilation if building on non-Linux
# -ldflags="-w -s" to make the binary smaller
RUN CGO_ENABLED=0 GOOS=linux go build -v -ldflags="-w -s" -o go-app .

# Stage 2: Create a lightweight image for the Go application
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the built executable from the builder stage & the templates
COPY --from=builder /app/go-app .
COPY --from=builder /app/templates ./templates

# Expose port 8081 (or whatever port your Go app listens on)
EXPOSE 8080

# Command to run the executable
# Database connection details (DB_HOST, DB_USER, etc.) will be passed in via docker-compose
CMD ["./go-app"]
