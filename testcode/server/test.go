package main

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var broker = "tcp://localhost:1883" // Đổi nếu cần
var username = "thaily"
var password = "Th@i2004"

func test() {
	// Cấu hình MQTT client
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("go_mqtt_client")

	// Thêm username/password nếu broker yêu cầu xác thực
	opts.SetUsername(username)
	opts.SetPassword(password)

	// Bật Auto-Reconnect (Tự động kết nối lại nếu bị mất kết nối)
	opts.SetAutoReconnect(true)

	// Bật Keep Alive (Giữ kết nối)
	opts.SetKeepAlive(60 * time.Second)

	// Xử lý khi mất kết nối
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		fmt.Println("Reconnected to MQTT broker!")
	})

	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		fmt.Println("Connection lost. Reconnecting...")
	})

	// Hàm xử lý khi nhận tin nhắn từ MQTT
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	})

	// Kết nối MQTT
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error()) // Kiểm tra lỗi kết nối
	}
	fmt.Println("Connected to MQTT broker!")

	// Đăng ký nhận dữ liệu từ topic "test/topic"
	topic := "test/topic"
	client.Subscribe(topic, 1, nil)
	fmt.Printf("Subscribed to topic: %s\n", topic)
	payload := "Hello from Go MQTT!"
	client.Publish(topic, 1, false, payload)
	fmt.Println("Published message:", payload)
	// Gửi tin nhắn liên tục lên MQTT
	for {
		time.Sleep(100 * time.Second) // Gửi tin nhắn mỗi 5 giây
	}
}
