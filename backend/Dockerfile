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

# Copy go.mod first (no go.sum — we'll generate it)
COPY go.mod ./

# Copy ALL source code so go mod tidy can see all import statements
COPY . .

# Run go mod tidy to:
#   1. Analyze all imports from source files
#   2. Generate go.sum automatically
#   3. Download all required module dependencies
# NOTE: go.sum is excluded from the build context via .dockerignore
#       so this step always generates a fresh, correct go.sum
RUN go mod tidy

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

# Copy migration and seed SQL files (read at runtime by the server)
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/seeds ./seeds

# Use non-root user
USER appuser

EXPOSE 8080

CMD ["./server"]
