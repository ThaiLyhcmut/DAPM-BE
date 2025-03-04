package controller

import (
	"ThaiLy/graph/model"
	protoKafka "ThaiLy/proto/kafka"
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/segmentio/kafka-go"
)

func (C *Controller) DeviceService(ctx context.Context, id int32, turnOn bool) (*string, error) {
	in := &protoKafka.DeviceRequest{
		Id:     id,
		TurnOn: turnOn,
	}
	res, err := C.kafka.DeviceService(ctx, in)
	if err != nil {
		return nil, err
	}
	return &res.Message, nil
}

func (C *Controller) DeviceStatusUpdated(ctx context.Context) (<-chan *model.Device, error) {
	ch := make(chan *model.Device)

	go func() {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"localhost:9092"},
			Topic:   "device_status",
			GroupID: "graphql-consumer",
		})
		defer reader.Close()

		for {
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Println("Kafka read error:", err)
				close(ch)
				return
			}

			log.Println("Received Kafka event:", string(msg.Value))
			trunOn := msg.Value[len(msg.Value)-1] == '1' // Giả định '1' là bật, '0' là tắt
			idStr := string(msg.Key)                     // Chuyển []byte -> string
			deviceID, err := strconv.Atoi(idStr)         // Chuyển string -> int
			if err != nil {
				fmt.Println("Lỗi chuyển đổi ID:", err)
				return
			}
			id := int32(deviceID) // Chuyển int -> int32
			ch <- &model.Device{ID: id, TurnOn: trunOn}
		}
	}()

	return ch, nil
}
