package controller

import (
	"ThaiLy/graph/config"
	"ThaiLy/graph/helper"
	"ThaiLy/graph/model"
	protoKafka "ThaiLy/proto/kafka"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/segmentio/kafka-go"
)

func (C *Controller) DeviceService(ctx context.Context, id int32, turnOn bool) (*string, error) {
	Claims, ok := ctx.Value(helper.Auth).(*helper.Claims)
	fmt.Print(Claims.ID, ok)
	if !ok {
		return nil, fmt.Errorf("Unauthorzition")
	}
	in := &protoKafka.DeviceRequest{
		Id:        id,
		TurnOn:    turnOn,
		AccountId: Claims.ID,
	}
	res, err := C.kafka.DeviceService(ctx, in)
	if err != nil {
		return nil, err
	}
	return &res.Message, nil
}

func (C *Controller) DeviceStatusUpdated(ctx context.Context) (<-chan *model.Device, error) {
	ch := make(chan *model.Device) // Tạo channel để gửi dữ liệu cho WebSocket client

	go func() {
		defer close(ch) // Đóng channel khi goroutine kết thúc

		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"kafka:9092"},
			Topic:   "device_status",
			GroupID: "graphql-consumer",
		})
		defer reader.Close()

		for {
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Println("Kafka read error:", err)
				return
			}

			log.Println("Received Kafka event:", string(msg.Value))

			// Tách dữ liệu từ message: "deviceId|turnOn|accountId"
			parts := strings.Split(string(msg.Value), "|")
			if len(parts) != 3 {
				log.Println("Invalid message format:", msg.Value)
				continue
			}

			// Chuyển đổi deviceId
			deviceID, err := strconv.Atoi(parts[0])
			if err != nil {
				log.Println("Invalid device ID:", parts[0])
				continue
			}

			// Chuyển đổi trạng thái turnOn
			turnOn := parts[1] == "true"

			// Chuyển đổi accountId
			accountId, err := strconv.Atoi(parts[2])
			if err != nil {
				log.Println("Invalid account ID:", parts[2])
				continue
			}

			// Kiểm tra accountId có đang kết nối không (trong MongoDB)
			isConnected, err := config.IsUserConnected(int32(accountId))
			if err != nil {
				log.Println("Error checking user connection:", err)
				continue
			}

			if !isConnected {
				log.Println("User ID", accountId, "is not connected.")
				continue
			}

			// Gửi WebSocket chỉ cho userId đang kết nối
			device := &model.Device{ID: int32(deviceID), TurnOn: turnOn}

			C.mu.Lock()
			if clientChan, exists := C.subscriptions[int32(accountId)]; exists {
				clientChan <- device
			} else {
				log.Println("No WebSocket subscription found for account ID:", accountId)
			}
			C.mu.Unlock()
		}
	}()

	return ch, nil
}
