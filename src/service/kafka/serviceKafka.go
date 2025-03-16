package main

import (
	protoKafka "ThaiLy/proto/kafka"
	"ThaiLy/service/kafka/controller"
	"ThaiLy/service/kafka/database"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

type DeviceService struct {
	c *controller.Controller
	protoKafka.UnimplementedDeviceServiceServer
}

func (s *DeviceService) ToggleDevice(ctx context.Context, req *protoKafka.DeviceRequest) (*protoKafka.DeviceResponse, error) {
	topic := "device_status"
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   topic,
	})
	defer writer.Close()

	// Gửi deviceId | turnOn | accountId
	message := fmt.Sprintf("%d|%t|%d", req.Id, req.TurnOn, req.AccountId)
	err := writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", req.Id)),
		Value: []byte(message),
	})
	if err != nil {
		return nil, err
	}

	return &protoKafka.DeviceResponse{Message: "Device state updated"}, nil
}

func handleMQTTMessage(client mqtt.Client, msg mqtt.Message) {
	kafkaTopic := "device_status"
	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"kafka:9092"},
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
	message := fmt.Sprintf("%d|%t|%d", req.Id, req.TurnOn, req.AccountId)

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

func connectMQTT(broker, clientID, username, password, topic string) {
	fmt.Println(broker, clientID, username, password, topic)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)

	// Thêm username và password
	opts.SetUsername(username)
	opts.SetPassword(password)

	opts.SetDefaultPublishHandler(handleMQTTMessage)

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

func main() {
	godotenv.Load()
	db := database.InitDB()
	// Ensure proper cleanup
	ctrl := controller.NewController(db)
	go connectMQTT(os.Getenv("MQTT_BROKER"), "myClient", os.Getenv("MQTT_USER"), os.Getenv("MQTT_PASSWORD"), os.Getenv("MQTT_TOPIC"))
	lis, err := net.Listen("tcp", "0.0.0.0:55557") // tao port
	if err != nil {
		log.Fatalf("err while create listen %v", err)
	}

	s := grpc.NewServer() // tao server

	protoKafka.RegisterDeviceServiceServer(s, &DeviceService{c: ctrl}) // dang ky server

	err = s.Serve(lis) // run server                                     // run server
	if err != nil {
		log.Fatalf("err while server %v", err)
	}
}
