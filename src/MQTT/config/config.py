import os

class Config:
    MQTT_BROKER = os.getenv("MQTT_BROKER")
    CLIENT_ID = os.getenv("CLIENT_ID")
    MQTT_USER = os.getenv("MQTT_USER")
    MQTT_PASSWORD = os.getenv("MQTT_PASSWORD")
    MQTT_TOPIC = os.getenv("MQTT_TOPIC")
    MY_HOST = os.getenv("MY_HOST")
    MY_PORT = os.getenv("MY_PORT")