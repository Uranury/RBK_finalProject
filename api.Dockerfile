FROM golang:1.23.1-alpine AS builder

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

COPY .env ./

# Build the API binary
RUN go build -o api cmd/api/*.go

# Production stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/.env ./

# Copy the binary
COPY --from=builder /app/api .

# Copy migrations if needed
COPY --from=builder /app/migrations ./migrations

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup
USER appuser

EXPOSE 8080

CMD ["./api"]