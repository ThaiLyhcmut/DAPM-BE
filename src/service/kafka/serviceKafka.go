package main

import (
	protoKafka "ThaiLy/proto/kafka"
	"ThaiLy/service/kafka/controller"
	"ThaiLy/service/kafka/database"
	"context"
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

type DeviceService struct {
	c *controller.Controller
	protoKafka.UnimplementedDeviceServiceServer
}

func (s *DeviceService) ToggleDevice(ctx context.Context, req *protoKafka.DeviceRequest) (*protoKafka.DeviceResponse, error) {
	// Gửi sự kiện trạng thái thiết bị vào Kafka
	topic := "device_status"
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   topic,
	})
	defer writer.Close()

	message := fmt.Sprintf("Device %d is now %t", req.Id, req.TurnOn)
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
