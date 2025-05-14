import os
import logging
import paho.mqtt.client as mqtt
import requests
import json
from controllers.controller import MQTTController  # Import class mới
client = mqtt.Client()
import asyncio

class MQTTConfig:
  def __init__(self):
    self.broker = os.getenv("MQTT_BROKER")
    self.port = int(os.getenv("MQTT_PORT"))
    self.client_id = os.getenv("CLIENT_ID")
    self.username = os.getenv("MQTT_USER")
    self.password = os.getenv("MQTT_PASSWORD")
    self.topics = [os.getenv("MQTT_TOPIC_DEVICE"), os.getenv("MQTT_TOPIC_AUDIO")]
    self.topics = list(set(self.topics))
    self.client = client
    self.aio_url = "https://io.adafruit.com/api/v2"
    self.controller = MQTTController(self.client)
    self.client.user_data_set({"status": "divices"})

  def connectMQTT(self):
    self.client.client_id = self.client_id
    self.client.username_pw_set(self.username, self.password)
    self.topics = self.controller.get_topic() # get topic on server
    try:
      self.client.connect(self.broker, self.port)
      logging.info("✅ Kết nối MQTT thành công!")

    except Exception as e:
      logging.error(f"❌ Lỗi kết nối MQTT: {e}")
      print("connect-error")
      return
    for topic in self.topics:
      self.client.subscribe(topic)
      self.create_topic(topic)    

    self.client.on_message = self.controller.controller
    self.client.loop_start()

  
  async def gql_listener(self):
    async for data in self.controller.socket_gql():
        device_id = data["id"]
        turn_on = data["turnOn"]
        if device_id is None:
            continue  # Bỏ qua nếu lỗi parse

        # Tìm name_id từ self.topics
        name_id = None
        for topic in self.topics:
            if topic.split("_")[-1] == str(device_id):
                name_id = topic
                break

        if name_id is None:
            logging.warning("device_id %s không tìm thấy trong self.topics", device_id)
            continue

        value = "ON" if turn_on else "OFF"
        self.client.user_data_set({"status": "server"})
        self.client.publish(name_id, value)
        logging.info("📡 Đã gửi lệnh tới thiết bị (%s): %s", name_id, value)

  def create_adafruit_feed(self, feed_name):
    """Tạo feed mới trên Adafruit IO nếu chưa tồn tại"""
    headers = {
        'X-AIO-Key': self.password,
        'Content-Type': 'application/json'
    }
    
    # Kiểm tra feed tồn tại
    check_url = f"{self.aio_url}/{self.username}/feeds/{feed_name}"
    try:
        response = requests.get(check_url, headers=headers)
        if response.status_code == 404:
            # Tạo feed mới nếu không tìm thấy
            create_url = f"{self.aio_url}/{self.username}/feeds"
            data = {
                'name': feed_name,
                'key': feed_name
            }
            response = requests.post(
                create_url, 
                headers=headers,
                data=json.dumps(data)
            )
            if response.status_code == 201:
                logging.info(f"✅ Đã tạo feed mới: {feed_name}")
                return True
            else:
                logging.error(f"❌ Không thể tạo feed: {response.text}")
                return False
        return True
    except Exception as e:
        logging.error(f"❌ Lỗi khi tạo feed: {e}")
        return False

  def create_topic(self, topic_name):
    """Tạo một topic mới và tự động tạo feed trên Adafruit nếu cần"""
    try:
      # Tách lấy tên feed từ topic (nếu topic có format username/feeds/feed-name)
      feed_name = topic_name.split('/')[-1] if '/' in topic_name else topic_name
      
      # Tạo feed trên Adafruit IO trước
      if self.create_adafruit_feed(feed_name):
          # Sau khi tạo feed thành công, thêm vào danh sách topics
          full_topic = f"{self.username}/feeds/{feed_name}"
          if full_topic not in self.topics:
              self.topics.append(full_topic)
              self.client.subscribe(full_topic)
              logging.info(f"✅ Đã tạo và subscribe topic: {full_topic}")
              return True
      return False
    except Exception as e:
      logging.error(f"❌ Lỗi khi tạo topic: {e}")
      return False

  def remove_topic(self, topic_name):
    """Xóa một topic khỏi danh sách theo dõi"""
    try:
      if topic_name in self.topics:
        self.topics.remove(topic_name)
        self.client.unsubscribe(topic_name)
        logging.info(f"✅ Đã xóa topic: {topic_name}")
        return True
      return False
    except Exception as e:
      logging.error(f"❌ Lỗi khi xóa topic: {e}")
      return False