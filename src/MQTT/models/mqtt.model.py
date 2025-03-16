import json
import logging
from paho.mqtt.client import Client as MQTTClient
from confluent_kafka import Producer
import os

# Khởi tạo Kafka producer
def create_kafka_producer():
    kafka_broker = os.getenv("KAFKA_BROKER")
    conf = {'bootstrap.servers': kafka_broker}
    producer = Producer(conf)
    return producer

# Xử lý message từ MQTT và ghi vào Kafka
def handle_mqtt_message(client, userdata, msg):
    kafka_topic = os.getenv("DEVICE_TOGGLE_TOPIC")

    # Parse nội dung MQTT message
    try:
        req = json.loads(msg.payload.decode())
    except json.JSONDecodeError as e:
        logging.error(f"❌ Lỗi parse MQTT message: {e}")
        return

    # Format lại dữ liệu giống ToggleDevice
    message = f"{req['Id']}|{req['TurnOn']}|{req['AccountId']}"

    # Ghi vào Kafka
    producer = create_kafka_producer()
    try:
        producer.produce(kafka_topic, key=str(req['Id']), value=message)
        producer.flush()
        logging.info(f"✅ Ghi vào Kafka thành công: {message}")
    except Exception as e:
        logging.error(f"❌ Lỗi ghi vào Kafka: {e}")
