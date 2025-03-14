# Stage 1: Build ứng dụng
FROM golang:1.24.0 AS builder

WORKDIR /app

# Kiểm tra file trước khi copy
RUN ls -lah 

# Copy toàn bộ project
COPY . .

# Cài đặt dependency
RUN go mod download

# Copy toàn bộ code vào container
COPY . .

# Build server
RUN go build -o server ./server/server.go

# Stage 2: Tạo image nhỏ gọn để chạy
FROM golang:1.23.0

WORKDIR /app

# Copy binary từ builder
COPY --from=builder /app/server .

# Mở port 8080
EXPOSE 8081

# Chạy server
CMD ["./server"]
