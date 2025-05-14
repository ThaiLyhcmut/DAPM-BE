import json
import logging
import asyncio
import sys
import os
from graphql_custom.graphql import GraphQL as LocalGraphQL

class MQTTController:
  def __init__(self, client):
    self.client = client
    self.graphql = LocalGraphQL()  # Use your local GraphQL class
  def get_topic(self):
    return self.graphql.get_home()
  async def socket_gql(self):
    async for msg in self.graphql.subscribe(self.client):
      yield msg  # Không dùng return
  def convert_to_json(self, msg):
    """Chuyển payload từ MQTT message thành JSON và trích xuất thông tin"""
    try:
        req = json.loads(msg.payload.decode())
    except json.JSONDecodeError as e:
        logging.error(f"❌ Lỗi parse MQTT message: {e}")
        return None, None, None  # Return three None values to match expected return structure
    
    message = f"{req['Id']}|{req['TurnOn']}|{req['AccountId']}"
    return int(req['Id']), bool(req['TurnOn']), message

  def controller(self, client, userdata, msg):
    print(userdata)
    """Xử lý tin nhắn nhận từ MQTT"""
    if userdata['status'] == 'server':
       userdata['status'] = "devices"
       return
    print("📩 Received message:", msg.payload.decode())
    device_id = msg.topic.split("_")[-1]
    if msg.topic == "audio":
      self.audio_server(msg, device_id)
    else:
      self.device_server(msg, device_id)

  def device_server(self, msg, device_id):
    """Xử lý tin nhắn liên quan đến thiết bị"""
    turnOn = (msg.payload.decode() == "ON")
    if device_id:
      print(f"✅ Processed Device ID: {device_id} - Message: {turnOn}")
      print(self.graphql.toggleDevice(device_id, turnOn))
    else:
      print("⚠ Lỗi xử lý MQTT message")
  
  def audio_server(self, msg, device_id):
    """Xử lý tin nhắn liên quan đến âm thanh"""
    device_id, turnOn, message = self.convert_to_json(msg)  # Fixed to match return values
    if device_id:
      print(f"✅ Processed Audio ID: {device_id} - Message: {message}")
      # self.process_audio(audio_id)  # This function is not defined in the original code
    else:
      print("⚠ Lỗi xử lý MQTT message")