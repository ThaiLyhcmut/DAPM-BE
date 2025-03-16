# Sử dụng Python 3.10 làm base image
FROM python:3.10

# Đặt thư mục làm việc trong container
WORKDIR /app

# Sao chép file mã nguồn vào container
COPY . .

# Cài đặt các thư viện cần thiết
RUN pip install --no-cache-dir paho-mqtt confluent-kafka

RUN python build -o iot ./IoT/main.py

# Chạy script Python
CMD ["./iot"]
