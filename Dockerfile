# Builder stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Runtime stage
FROM alpine:3.19
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Expose
EXPOSE 8080

# Run
CMD ["./server"]
