FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

RUN go build -o worker cmd/worker/*.go

# Production stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/worker .

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup
USER appuser

CMD ["./worker"]