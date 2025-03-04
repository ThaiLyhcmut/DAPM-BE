package controller

import "ThaiLy/server/client"

type Controller struct {
	auth      *client.GRPCAuthClient
	equipment *client.GRPCEquipmentClient
}
