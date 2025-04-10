package controller

import (
	"ThaiLy/graph/config"
	"ThaiLy/graph/helper"
	"ThaiLy/graph/model"
	protoKafka "ThaiLy/proto/kafka"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/segmentio/kafka-go"
)

func (C *Controller) DeviceService(ctx context.Context, id int32, turnOn bool) (*string, error) {
	Claims, ok := ctx.Value(helper.Auth).(*helper.Claims)
	if !ok {
		return nil, fmt.Errorf("Unauthorzition")
	}
	Id, err := helper.ParseASE(Claims.ID)

	if err != nil {
		return nil, fmt.Errorf("error parse id")
	}

	equipment, err := C.equipment.CheckEquipment(id)

	if err != nil {
		return nil, fmt.Errorf("error get Equipment by id")
	}

	checkEquipmentByAccountId, err := C.equipment.CheckHome(Id, equipment.HomeId)

	if err != nil && checkEquipmentByAccountId == nil {
		return nil, fmt.Errorf("error equipment id not your home")
	}

	if _, err = C.equipment.ChangeTurnOnEquipment(id, turnOn); err != nil {
		return nil, err
	}

	in := &protoKafka.DeviceRequest{
		Id:        id,
		TurnOn:    turnOn,
		AccountId: Id,
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
	Id, err := helper.ParseASE(Claims.ID)
	if err != nil {
		return nil, fmt.Errorf("error parse id")
	}
	// Lưu channel vào danh sách
	userChannels[Id] = append(userChannels[Id], ch)

	go func() {
		defer func() {
			// Xóa channel khi đóng
			for i, c := range userChannels[Id] {
				if c == ch {
					userChannels[Id] = append(userChannels[Id][:i], userChannels[Id][i+1:]...)
					break
				}
			}
			close(ch)
		}()

		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{os.Getenv("KAFKA_BROKER")},
			Topic:   os.Getenv("DEVICE_TOGGLE_TOPIC"),
			GroupID: os.Getenv("GROUP_ID"),
		})

		for {
			msg, err := reader.ReadMessage(ctx)

			if err != nil {
				log.Println("Kafka read error:", err)
				return
			}

			log.Println("Received Kafka event: ", string(msg.Value))
			parts := strings.Split(string(msg.Value), "|")
			if len(parts) != 3 {
				log.Println("Invalid message format:", msg.Value)
				continue
			}

			deviceID, err := strconv.Atoi(parts[0])
			turnOn := parts[1] == "true"
			accountId, err := strconv.ParseInt(string(parts[2]), 10, 32)
			if err != nil {
				return
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
