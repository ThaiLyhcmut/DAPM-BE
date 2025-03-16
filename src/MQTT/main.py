import logging
from flask import Flask
from controllers.mqtt_controller import start_mqtt
from views.mqtt_view import home
from threading import Thread
import os
from dotenv import load_dotenv
load_dotenv()
# Tạo ứng dụng Flask
app = Flask(__name__)

# Đăng ký route
app.add_url_rule('/', 'home', home)

# Khởi động MQTT trong một thread riêng
def start_mqtt_thread():
    broker = os.getenv("MQTT_BROKER")
    client_id = os.getenv("CLIENT_ID")
    username = os.getenv("MQTT_USER")
    password = os.getenv("MQTT_PASSWORD")
    topic = os.getenv("MQTT_TOPIC")
    port = int(os.getenv("MQTT_PORT"))
    start_mqtt(broker, port, client_id, username, password, topic)

if __name__ == "__main__":
    mqtt_thread = Thread(target=start_mqtt_thread)
    mqtt_thread.daemon = True  # Đảm bảo thread này sẽ tự động tắt khi Flask tắt
    mqtt_thread.start()
    app.run(host=os.getenv("MY_HOST"), port=os.getenv("MY_PORT"))
