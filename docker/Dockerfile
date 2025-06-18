# =============================================
# Build Stage (Multi-purpose with CGO toggle)
# =============================================

# Use Alpine-based Go image for small size (~150MB)
# Switch to 'golang:1.24-bullseye' if CGO is needed
FROM golang:1.24-alpine AS builder

# Set working directory for all subsequent commands
WORKDIR /app

# Copy dependency files first to leverage Docker cache
# Note: Ensure go.mod exists in build context
COPY ../go.mod ../go.sum ./
RUN go mod download && go mod verify

# Copy all application source code
# Note: Use .dockerignore to exclude unnecessary files
COPY .. .

# Build configuration:
# - CGO_ENABLED=0: Static binary (default) - comment out for CGO
# - GOOS=linux: Ensure Linux compatibility
# - -ldflags="-s -w": Strip debug symbols (-s) and DWARF (-w)
# - -o /app/main: Output path for compiled binary
# - ./...: Build all packages in directory
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/main ./.

# =============================================
# Runtime Stage (Secure minimal deployment)
# =============================================

# Use Alpine for small size (~5MB) with basic shell
# Alternatives:
# - FROM scratch: For absolute minimal (no shell)
# - FROM gcr.io/distroless/static: Google's secure minimal image
FROM alpine:3.20

# Set working directory for runtime
WORKDIR /app

# Create non-root user/group for security
# - -S: System user (no password, no home dir)
# - -G: Assign to group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Transfer ownership of /app to appuser
RUN chown -R appuser:appgroup /app

# Switch to non-root user context
USER appuser

# Copy only the compiled binary from builder
COPY --from=builder --chown=appuser:appgroup /app/main .

# Expose application port
EXPOSE 8080

# Runtime command
CMD ["./main"]