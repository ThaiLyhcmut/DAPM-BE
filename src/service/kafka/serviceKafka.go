package main

import (
	protoKafka "ThaiLy/proto/kafka"
	"ThaiLy/service/kafka/controller"
	"ThaiLy/service/kafka/database"
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

type DeviceService struct {
	c *controller.Controller
	protoKafka.UnimplementedDeviceServiceServer
}

func (s *DeviceService) ToggleDevice(ctx context.Context, req *protoKafka.DeviceRequest) (*protoKafka.DeviceResponse, error) {
	topic := os.Getenv("DEVICE_TOGGLE_TOPIC")
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKER")},
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

func main() {
	godotenv.Load()
	db := database.InitDB()
	// Ensure proper cleanup
	ctrl := controller.NewController(db)
	lis, err := net.Listen(os.Getenv("NET_WORK"), os.Getenv("ADDRESS")) // tao port
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
