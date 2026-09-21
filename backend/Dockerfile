# ==============================================================
# MULTI-STAGE DOCKERFILE — Go Backend
# Stage 1 (builder): Compile Go binary
# Stage 2 (runtime): Minimal Alpine image (~20MB total)
# ==============================================================

# ---- Stage 1: Build ----
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy dependency manifests first for layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Build a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /app/server ./cmd/server/main.go

# ---- Stage 2: Runtime ----
FROM alpine:3.19

# Add CA certificates (needed for HTTPS calls if any)
RUN apk add --no-cache ca-certificates tzdata

# Create a non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy only the compiled binary from builder stage
COPY --from=builder /app/server .

# Copy migration and seed SQL files (embedded at runtime)
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/seeds ./seeds

# Use non-root user
USER appuser

EXPOSE 8080

CMD ["./server"]

