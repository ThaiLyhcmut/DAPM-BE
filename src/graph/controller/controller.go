package controller

import "ThaiLy/server/client"

type Controller struct {
	auth      *client.GRPCAuthClient
	equipment *client.GRPCEquipmentClient
}

// Constructor function
func NewController(auth *client.GRPCAuthClient, equipment *client.GRPCEquipmentClient) *Controller {
	return &Controller{auth: auth, equipment: equipment}
}
