import logging
from flask import Flask
from config.config import MQTTConfig
from views.mqtt_view import home
from threading import Thread
import os
from dotenv import load_dotenv
import asyncio
import time
load_dotenv(".mqtt.env")

# Tạo ứng dụng Flask
app = Flask(__name__)

# Đăng ký route
app.add_url_rule('/', 'home', home)

mqtt_config = MQTTConfig()

# Hàm chạy asyncio event loop trong thread riêng
def run_async_loop(loop):
    asyncio.set_event_loop(loop)
    loop.run_forever()

# Add this function to monitor terminal input
def monitor_terminal_input(mqtt_config):
    while True:
        try:
            command = input().strip()
            if command.lower() == "reload":
                print("🔄 Reloading MQTT connection...")
                # Disconnect first if connected
                try:
                    mqtt_config.client.disconnect()
                    mqtt_config.client.loop_stop()
                    print("Disconnected from previous MQTT session")
                except:
                    pass
                
                # Reconnect
                time.sleep(1)  # Small delay to ensure clean disconnect
                mqtt_config.connectMQTT()
                print("✅ MQTT connection reloaded successfully!")
        except KeyboardInterrupt:
            break
        except Exception as e:
            print(f"Error processing command: {e}")

# Khởi động MQTT trong một thread riêng
def start_mqtt_thread():
    mqtt_config.connectMQTT()

# Hàm khởi động gql_listener trong asyncio loop
async def start_gql_listener():
    await mqtt_config.gql_listener()
if __name__ == "__main__":
    # Tạo event loop mới cho asyncio
    loop = asyncio.new_event_loop()
    
    # Tạo và khởi động thread cho asyncio loop
    asyncio_thread = Thread(target=run_async_loop, args=(loop,))
    asyncio_thread.daemon = True
    asyncio_thread.start()
    
    # Lên lịch chạy gql_listener trong event loop
    asyncio.run_coroutine_threadsafe(start_gql_listener(), loop)
    
    # Khởi động thread MQTT
    mqtt_thread = Thread(target=start_mqtt_thread)
    mqtt_thread.daemon = True
    mqtt_thread.start()
    
    # Add terminal monitoring thread
    terminal_thread = Thread(target=monitor_terminal_input, args=(mqtt_config,))
    terminal_thread.daemon = True
    terminal_thread.start()
    
    # Khởi động Flask
    app.run(host=os.getenv("MY_HOST"), port=os.getenv("MY_PORT"))