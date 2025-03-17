import logging
from flask import Flask
from config.config import MQTTConfig
from views.mqtt_view import home
# from controllers.mqtt_controller import start_device
from threading import Thread
import os
from dotenv import load_dotenv

load_dotenv(".mqtt.env")
# Tạo ứng dụng Flask
app = Flask(__name__)

# Đăng ký route
app.add_url_rule('/', 'home', home)

# Khởi động MQTT trong một thread riêng
def start_mqtt_thread():
  mqtt_config = MQTTConfig()
  mqtt_config.connectMQTT()
  # start_device()

if __name__ == "__main__":
    mqtt_thread = Thread(target=start_mqtt_thread)
    mqtt_thread.daemon = True  # Đảm bảo thread này sẽ tự động tắt khi Flask tắt
    mqtt_thread.start()
    app.run(host=os.getenv("MY_HOST"), port=os.getenv("MY_PORT"))
