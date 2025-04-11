# Stage 1: Build Frontend
FROM node:20-slim AS frontend-builder

# Set working directory
WORKDIR /app/front

# Copy frontend source files
COPY front .

# Install dependencies and build
RUN npm install && npm run build

# Stage 2: Build Backend
FROM golang:1.23 AS backend-builder

# Set working directory
WORKDIR /app

# Copy backend source files
COPY . .

# Copy built frontend files
COPY --from=frontend-builder /app/public ./public

# Build backend
RUN go build -o tiny-nav

# Stage 3: Final Image
FROM alpine:latest

# Set environment variables
ENV LISTEN_PORT=58080

# Set working directory
WORKDIR /app

# Copy built binary and frontend files
COPY --from=backend-builder /app/tiny-nav .

# Expose port
EXPOSE 58080

# Run the application
CMD ["./tiny-nav"]
