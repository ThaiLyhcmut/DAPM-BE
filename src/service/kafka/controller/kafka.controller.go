package controller

import "ThaiLy/service/kafka/database"

type Controller struct {
	d *database.Database
}

func NewController(db *database.Database) *Controller {
	return &Controller{d: db}
}
