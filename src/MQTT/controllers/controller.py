import json
import logging
import graphql.graphql as gql

class MQTTController:
  def __init__(self, client):
    self.client = client
    self.graphql = gql.graphql()

  def convert_to_json(self, msg):
    """Chuyển payload từ MQTT message thành JSON và trích xuất thông tin"""
    try:
        req = json.loads(msg.payload.decode())
    except json.JSONDecodeError as e:
        logging.error(f"❌ Lỗi parse MQTT message: {e}")
        return None, None
    
    message = f"{req['Id']}|{req['TurnOn']}|{req['AccountId']}"
    return int(req['Id']), bool(req['TurnOn']), message

  def controller(self, client, userdata, msg):
    """Xử lý tin nhắn nhận từ MQTT"""
    print("📩 Received message:", msg.payload.decode())
    print(f"Topic: {msg.topic}, QoS: {msg.qos}, Retain: {msg.retain}")
    print("Client:", client)
    print("Message:", msg)

    if msg.topic == "audio":
      self.audio_server(msg)
    else:
      self.device_server(msg)

  def device_server(self, msg):
    """Xử lý tin nhắn liên quan đến thiết bị"""
    device_id, turnOn, message = self.convert_to_json(msg)
    if device_id:
      print(f"✅ Processed Device ID: {device_id} - Message: {message}")
      print(self.graphql.toggleDevice(device_id, turnOn))
    else:
      print("⚠ Lỗi xử lý MQTT message")
    # print(self.graphql.toggleDevice(device_id, turnOn))
  
  def audio_server(self, msg):
    """Xử lý tin nhắn liên quan đến âm thanh"""
    audio_id, message = self.convert_to_json(msg)
    if audio_id:
      print(f"✅ Processed Audio ID: {audio_id} - Message: {message}")
      self.process_audio(audio_id)
    else:
      print("⚠ Lỗi xử lý MQTT message")
