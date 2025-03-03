# Stage 1: Build ứng dụng
FROM golang:1.24.0 AS builder

WORKDIR /app

# Copy toàn bộ project
COPY . .

# Cài đặt dependency
RUN go mod download

# Build service-auth
RUN go build -o serviceAuth ./service/auth/serviceAuth.go

# Stage 2: Tạo image nhỏ gọn để chạy
FROM golang:1.23.0

WORKDIR /app

# Copy binary từ builder
COPY --from=builder /app/serviceAuth .

# Mở port 50051 cho gRPC
EXPOSE 55555

# Chạy service-auth
CMD ["./serviceAuth"]
