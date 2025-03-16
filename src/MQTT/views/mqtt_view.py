from flask import jsonify

def home():
    return jsonify({"message": "✅ MQTT Listener is running!"})
