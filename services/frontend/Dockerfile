# --------------------------
# Stage 1: Build the Vite client
# --------------------------
FROM node:20-alpine AS client-builder

WORKDIR /app

# Copy package files and install dependencies
COPY client/package*.json ./
COPY client/vite.config.ts ./
RUN npm ci

# Copy the rest of the client source code
COPY client/ ./

# Build the client
RUN npm run build

# --------------------------
# Stage 2: Build the Go server
# --------------------------
FROM golang:1.25-alpine AS go-builder

WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy Go source code
COPY ./ ./

# Copy built client
COPY --from=client-builder /app/dist ./client/dist

# Build the Go server binary
RUN go build -o server ./cmd/main.go

# --------------------------
# Stage 3: Final minimal image
# --------------------------
FROM alpine:3.18

WORKDIR /app

# Copy Go binary and client dist
COPY --from=go-builder /app/server ./server
COPY --from=go-builder /app/client/dist ./client/dist

# Optional: copy .env file if present
# COPY .env ./

# Expose port
EXPOSE 8080

# Run the server
CMD ["./server"]
