package controller

import (
	"ThaiLy/graph/model"
	"ThaiLy/server/client"
	"sync"
)

type Controller struct {
	auth          *client.GRPCAuthClient
	equipment     *client.GRPCEquipmentClient
	kafka         *client.GRPCKafkaClient
	mu            sync.Mutex
	subscriptions map[int32]chan *model.Device
}

// Constructor function
func NewController(auth *client.GRPCAuthClient, equipment *client.GRPCEquipmentClient, kafka *client.GRPCKafkaClient) *Controller {
	return &Controller{auth: auth, equipment: equipment, kafka: kafka}
}
