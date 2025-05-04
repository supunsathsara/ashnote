# Start from the official golang image
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ashnote ./cmd/main.go

# Use a small alpine image for the final container
FROM alpine:latest

# Install necessary dependencies for SQLite
RUN apk --no-cache add ca-certificates tzdata

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/ashnote .

# Copy web templates
COPY --from=builder /app/web/templates ./web/templates

# Expose port 3000
EXPOSE 3000

# Command to run the executable
CMD ["./ashnote"]