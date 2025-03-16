import os
import json
import paho.mqtt.client as mqtt
from confluent_kafka import Producer

# Lấy thông tin từ biến môi trường
MQTT_BROKER = os.getenv("MQTT_BROKER", "localhost")
MQTT_PORT = int(os.getenv("MQTT_PORT", "1883"))
MQTT_USERNAME = os.getenv("MQTT_USERNAME", "thaily")
MQTT_PASSWORD = os.getenv("MQTT_PASSWORD", "Th@i2004")
MQTT_TOPIC = os.getenv("MQTT_TOPIC", "test/topic")

print(MQTT_BROKER, MQTT_PORT, MQTT_USERNAME, MQTT_PASSWORD, MQTT_TOPIC)

KAFKA_BROKER = os.getenv("KAFKA_BROKER", "kafka:9092")
KAFKA_TOPIC = os.getenv("KAFKA_TOPIC", "device_status")

print(KAFKA_BROKER, KAFKA_TOPIC)
# Khởi tạo Kafka Producer
producer = Producer({"bootstrap.servers": KAFKA_BROKER})

# Callback khi kết nối MQTT thành công
def on_connect(client, userdata, flags, rc):
    if rc == 0:
        print("✅ Kết nối MQTT thành công!")
        client.subscribe(MQTT_TOPIC)
    else:
        print(f"❌ Kết nối MQTT thất bại với mã lỗi: {rc}")

# Callback khi nhận tin nhắn từ MQTT
def on_message(client, userdata, msg):
    message = msg.payload.decode()
    print(f"📩 Nhận tin nhắn: {message} từ topic {msg.topic}")
    try:
        data = json.loads(message)
        formatted_message = f"{data['Id']}|{data['TurnOn']}|{data['AccountId']}"

        # Gửi tin nhắn vào Kafka topic
        producer.produce(KAFKA_TOPIC, key=str(data['Id']), value=formatted_message.encode())
        producer.flush()
        print(f"🚀 Đã gửi tin nhắn vào Kafka topic: {KAFKA_TOPIC} với nội dung: {formatted_message}")

    except (json.JSONDecodeError, KeyError) as e:
        print(f"❌ Lỗi xử lý tin nhắn: {e}")

# Tạo MQTT Client
client = mqtt.Client()
client.username_pw_set(MQTT_USERNAME, MQTT_PASSWORD)

# Gán Callback
client.on_connect = on_connect
client.on_message = on_message

# Kết nối MQTT
client.connect(MQTT_BROKER, MQTT_PORT, 60)

# Chạy vòng lặp nhận dữ liệu
client.loop_forever()
