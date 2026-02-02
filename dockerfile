FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kafka-mongo-service .

# Use a minimal alpine image for the final stage
FROM alpine:3.17

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/kafka-mongo-service .

# Create a non-root user to run the application
RUN adduser -D -g '' appuser
USER appuser

# Command to run the executable
CMD ["./kafka-mongo-service"]