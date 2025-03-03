package main

import (
	"context"
	"log"
	"net"

	protoEquipment "ThaiLy/proto/equipment"
	"ThaiLy/service/equipment/controller"
	"ThaiLy/service/equipment/database"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type service struct {
	c *controller.Controller
	protoEquipment.UnimplementedEquipmentServiceServer
}

func (S *service) Home(ctx context.Context, in *protoEquipment.HomeRQ) (*protoEquipment.ListHomeRP, error) {
	return S.c.ControllerHome(in.GetAccountId())
}

func (S *service) Area(ctx context.Context, in *protoEquipment.AreaRQ) (*protoEquipment.ListAreaRP, error) {
	return S.c.ControllerArea(in.GetHomeId())
}

func (S *service) Equipment(ctx context.Context, in *protoEquipment.EquipmentRQ) (*protoEquipment.ListEquimentRP, error) {
	return S.c.ControllerEquipment(in.GetHomeId())
}

func main() {
	godotenv.Load()
	db := database.InitDB()
	// Ensure proper cleanup
	ctrl := controller.NewController(db)

	lis, err := net.Listen("tcp", "0.0.0.0:55556") // tao port
	if err != nil {
		log.Fatalf("err while create listen %v", err)
	}

	s := grpc.NewServer() // tao server

	protoEquipment.RegisterEquipmentServiceServer(s, &service{c: ctrl}) // dang ky server

	err = s.Serve(lis) // run server                                     // run server
	if err != nil {
		log.Fatalf("err while server %v", err)
	}
}
