# Sử dụng image Python chính thức
FROM python:3.12-slim

# Thiết lập thư mục làm việc
WORKDIR /app

# Sao chép các file cần thiết vào Docker container
COPY requirements.txt .

# Cài đặt các thư viện cần thiết
RUN pip install --no-cache-dir -r requirements.txt

# Sao chép toàn bộ mã nguồn vào container
COPY . .


# Thiết lập biến môi trường cho Flask
ENV FLASK_APP=MQTT/main.py
ENV FLASK_ENV=production

# Expose port cho Flask (5000)
EXPOSE 5000

# Chạy ứng dụng Flask
CMD ["python", "MQTT/main.py"]