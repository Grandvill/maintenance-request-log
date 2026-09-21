# ==============================================================
# MULTI-STAGE DOCKERFILE — Nuxt 3 Frontend
# Stage 1 (deps):    Install dependencies
# Stage 2 (builder): Build the Nuxt app
# Stage 3 (runtime): Minimal Node.js runtime image
# ==============================================================

# ---- Stage 1: Install Dependencies ----
# Use node:24-alpine to match the local development environment (node v24)
# This prevents npm lock file version mismatch errors
FROM node:24-alpine AS deps

WORKDIR /app

COPY package.json package-lock.json* ./

# Use npm install (not ci) to handle lock file regeneration gracefully
# This is intentional: lock file generated locally on Node 24 is used here
RUN npm install --no-audit --no-fund --prefer-offline 2>/dev/null || npm install --no-audit --no-fund

# ---- Stage 2: Build ----
FROM node:24-alpine AS builder

WORKDIR /app

# Accept build-time env variable for public API URL
ARG NUXT_PUBLIC_API_BASE_URL
ENV NUXT_PUBLIC_API_BASE_URL=${NUXT_PUBLIC_API_BASE_URL}

# Copy dependencies from deps stage
COPY --from=deps /app/node_modules ./node_modules
COPY . .

# Build the Nuxt application
RUN npm run build

# ---- Stage 3: Runtime ----
FROM node:24-alpine AS runtime

WORKDIR /app

# Create a non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy only the built output (standalone Nitro server)
COPY --from=builder /app/.output ./

# Use non-root user
USER appuser

EXPOSE 3000

ENV NITRO_PORT=3000
ENV NITRO_HOST=0.0.0.0

CMD ["node", "server/index.mjs"]
