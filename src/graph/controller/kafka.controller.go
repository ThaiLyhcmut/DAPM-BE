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
	IDP, err := helper.ParseASE(Claims.ID)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("error parse id")
	}
	in := &protoKafka.DeviceRequest{
		Id:        id,
		TurnOn:    turnOn,
		AccountId: IDP,
	}
	res, err := C.kafka.DeviceService(ctx, in)
	if err != nil {
		return nil, err
	}
	return &res.Message, nil
}

var userChannels = make(map[int32][]chan *model.Device)

func (C *Controller) DeviceStatusUpdated(ctx context.Context) (<-chan *model.Device, error) {
	ch := make(chan *model.Device)

	Claims, ok := ctx.Value(helper.Auth).(*helper.Claims)
	if !ok {
		return nil, fmt.Errorf("could not retrieve claims from context")
	}
	IDP, err := helper.ParseASE(Claims.ID)
	if err != nil {
		return nil, err
	}
	Id, err := strconv.Atoi(IDP)
	if err != nil {
		return nil, fmt.Errorf("error parse id")
	}
	userID := int32(Id)

	// Lưu channel vào danh sách
	userChannels[userID] = append(userChannels[userID], ch)

	go func() {
		defer func() {
			// Xóa channel khi đóng
			for i, c := range userChannels[userID] {
				if c == ch {
					userChannels[userID] = append(userChannels[userID][:i], userChannels[userID][i+1:]...)
					break
				}
			}
			close(ch)
		}()

		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"kafka:9092"},
			Topic:   "device_status",
			GroupID: "graphql-consumer",
		})

		for {
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				log.Println("Kafka read error:", err)
				return
			}

			log.Println("Received Kafka event:", string(msg.Value))
			parts := strings.Split(string(msg.Value), "|")
			if len(parts) != 3 {
				log.Println("Invalid message format:", msg.Value)
				continue
			}

			deviceID, err := strconv.Atoi(parts[0])
			turnOn := parts[1] == "true"
			accountId, err := strconv.Atoi(parts[2])
			if err != nil {
				log.Println("Invalid account ID:", parts[2])
				continue
			}
			device := &model.Device{ID: int32(deviceID), TurnOn: turnOn}
			for _, c := range userChannels[int32(accountId)] {
				c <- device
			}
		}
	}()

	return ch, nil
}

func HandleDeviceUpdates(userID int32, deviceChan <-chan *model.Device) {
	log.Printf("Cập nhật thiết bị cho user %d: ", userID)

	err := config.SaveUserConnection(userID)
	if err != nil {
		log.Printf("Lỗi cập nhật thiết bị cho user %d: %v", userID, err)
	}
}
