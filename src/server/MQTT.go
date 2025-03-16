package main

import (
	protoKafka "ThaiLy/proto/kafka"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/segmentio/kafka-go"
)

func HandleMQTTMessage(client mqtt.Client, msg mqtt.Message) {
	kafkaTopic := os.Getenv("DEVICE_TOGGLE_TOPIC")
	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKER")},
		Topic:   kafkaTopic,
	})
	defer kafkaWriter.Close()

	// Parse nội dung MQTT message
	var req protoKafka.DeviceRequest
	err := json.Unmarshal(msg.Payload(), &req)
	if err != nil {
		log.Printf("❌ Lỗi parse MQTT message: %v", err)
		return
	}

	// Format lại dữ liệu giống ToggleDevice
	message := fmt.Sprintf("%d|%t|%s", req.Id, req.TurnOn, req.AccountId)

	// Ghi vào Kafka
	err = kafkaWriter.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", req.Id)),
		Value: []byte(message),
	})
	if err != nil {
		log.Printf("❌ Lỗi ghi vào Kafka: %v", err)
	} else {
		log.Printf("✅ Ghi vào Kafka thành công: %s", message)
	}
}

func ConnectMQTT(broker, clientID, username, password, topic string) {
	fmt.Println(broker, clientID, username, password, topic)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)

	// Thêm username và password
	opts.SetUsername(username)
	opts.SetPassword(password)

	opts.SetDefaultPublishHandler(HandleMQTTMessage)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("❌ Lỗi kết nối MQTT: %v", token.Error())
	}
	log.Println("✅ Kết nối MQTT thành công!")

	// Đăng ký lắng nghe topic MQTT
	if token := client.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
		log.Fatalf("❌ Lỗi đăng ký topic MQTT: %v", token.Error())
	}
	log.Printf("📩 Đang lắng nghe MQTT topic: %s", topic)
}
