# Stage 1: Build ứng dụng
FROM golang:1.24.0 AS builder

WORKDIR /app

# Copy toàn bộ project
COPY ../go.mod ../go.sum ./

# Cài đặt dependency
RUN go mod download

# Build service-auth
RUN go build -o serviceEquipment ./service/equipment/serviceEquipment.go

# Stage 2: Tạo image nhỏ gọn để chạy
FROM golang:1.23.0

WORKDIR /app

# Copy binary từ builder
COPY --from=builder /app/serviceEquipment .

# Mở port 55556 cho gRPC
EXPOSE 55556

# Chạy service-auth
CMD ["./serviceEquipment"]
