package controller

import "ThaiLy/server/client"

type Controller struct {
	auth      *client.GRPCAuthClient
	equipment *client.GRPCEquipmentClient
	kafka     *client.GRPCKafkaClient
}

// Constructor function
func NewController(auth *client.GRPCAuthClient, equipment *client.GRPCEquipmentClient, kafka *client.GRPCKafkaClient) *Controller {
	return &Controller{auth: auth, equipment: equipment, kafka: kafka}
}
