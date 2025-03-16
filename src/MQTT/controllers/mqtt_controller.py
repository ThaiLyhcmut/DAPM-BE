from models.mqtt_model import handle_mqtt_message
from paho.mqtt.client import Client as MQTTClient
import logging

def start_mqtt(broker, port, client_id, username, password, topic):
    print(broker, port, client_id, username, password, topic)
    client = MQTTClient()
    client.client_id = client_id
    client.username_pw_set(username, password)
    client.on_message = handle_mqtt_message

    try:
        client.connect(broker, port)
        logging.info("✅ Kết nối MQTT thành công!")
    except Exception as e:
        logging.error(f"❌ Lỗi kết nối MQTT: {e}")
        return

    client.subscribe(topic)
    logging.info(f"📩 Đang lắng nghe MQTT topic: {topic}")
    client.loop_start()
