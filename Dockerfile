# --------------------------
# Stage 1: Build the Vite client
# --------------------------
FROM node:20-alpine AS client-builder

# Set working directory
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

# Set working directory
WORKDIR /app

# Copy Go modules files and download dependencies
COPY server/go.mod server/go.sum ./
RUN go mod download

# Copy the Go server source code
COPY ./server/ ./

# Copy the built client from the previous stage
COPY --from=client-builder /app/dist ./client/dist

# Build the Go server binary
RUN go build server .

# --------------------------
# Stage 3: Final minimal image
# --------------------------
FROM alpine:3.18

# Install CA certificates (for HTTPS, etc.)
# RUN apk add --no-cache ca-certificates

# Set working directory
WORKDIR /app

# Copy the Go server binary and client
COPY --from=go-builder /app/server ./
COPY --from=go-builder /app/client/dist ./client/dist

# Expose port
EXPOSE 8080

# Run the server
CMD ["./server"]