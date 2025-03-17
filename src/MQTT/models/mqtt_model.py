import json
import logging
from paho.mqtt.client import Client as MQTTClient
from confluent_kafka import Producer
import os

# Khởi tạo Kafka producer
def create_kafka_producer():
    kafka_broker = os.getenv("KAFKA_BROKER")
    client_id = os.getenv("CLIENT_ID")
    print(kafka_broker)
    conf = {'bootstrap.servers': kafka_broker, 'client.id': client_id}
    print(conf)
    producer = Producer(conf)
    return producer

# Xử lý message từ MQTT và ghi vào Kafka
def handle_mqtt_message(client, userdata, msg):
    kafka_topic = os.getenv("DEVICE_TOGGLE_TOPIC")
    print(1)
    # Parse nội dung MQTT message
    try:
        req = json.loads(msg.payload.decode())
        print(5)
    except json.JSONDecodeError as e:
        print(6)
        logging.error(f"❌ Lỗi parse MQTT message: {e}")
        return

    print(2)
    # Format lại dữ liệu giống ToggleDevice
    message = f"{req['Id']}|{req['TurnOn']}|{req['AccountId']}"
    print(3)
    # Ghi vào Kafka
    producer = create_kafka_producer()
    print(4)
    try:  
        print(7)
        print(req, kafka_topic, message)
        producer.produce(kafka_topic, key="hello", value="world")
        print(10)
        producer.flush()
        print(9)
        logging.info(f"✅ Ghi vào Kafka thành công: {message}")
        print(11)
    except Exception as e:
        print(8)
        logging.error(f"❌ Lỗi ghi vào Kafka: {e}")
    print(12)
