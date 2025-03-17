import os
import logging
import paho.mqtt.client as mqtt
from controllers.controller import MQTTController  # Import class mới

client = mqtt.Client()

class MQTTConfig:
  def __init__(self):
    self.broker = os.getenv("MQTT_BROKER")
    self.port = int(os.getenv("MQTT_PORT"))
    self.client_id = os.getenv("CLIENT_ID")
    self.username = os.getenv("MQTT_USER")
    self.password = os.getenv("MQTT_PASSWORD")
    self.topic_device = os.getenv("MQTT_TOPIC_DEVICE")
    self.topic_audio = os.getenv("MQTT_TOPIC_AUDIO")
  def connectMQTT(self):
    client.client_id = self.client_id
    client.username_pw_set(self.username, self.password)

    try:
      client.connect(self.broker, self.port)
      logging.info("✅ Kết nối MQTT thành công!")
    except Exception as e:
      logging.error(f"❌ Lỗi kết nối MQTT: {e}")
      return
    client.subscribe(self.topic_device)
    client.subscribe(self.topic_audio)
    mqtt_controller = MQTTController(client)
    client.on_message = mqtt_controller.controller 
    client.loop_start()