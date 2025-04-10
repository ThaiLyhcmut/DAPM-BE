package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

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

func (S *service) CreateHome(ctx context.Context, in *protoEquipment.CreateHomeRQ) (*protoEquipment.HomeRP, error) {
	return S.c.ControllerCreateHome(in.AccountId, in.HomeName, in.Location, in.Deleted, in.CreatedAt)
}

func (S *service) DeleteHome(ctx context.Context, in *protoEquipment.DeleteHomeRQ) (*protoEquipment.SuccessRP, error) {
	return S.c.ControllerDeleteHome(in.Id)
}

func (S *service) EditHome(ctx context.Context, in *protoEquipment.EditHomeRQ) (*protoEquipment.HomeRP, error) {
	return nil, nil
}

func (S *service) Area(ctx context.Context, in *protoEquipment.AreaRQ) (*protoEquipment.ListAreaRP, error) {
	return S.c.ControllerArea(in.GetHomeId())
}

func (S *service) CreateArea(ctx context.Context, in *protoEquipment.CreateAreaRQ) (*protoEquipment.AreaRP, error) {
	return S.c.ControllerCreateArea(in.HomeId, in.Name)
}

func (S *service) DeleteArea(ctx context.Context, in *protoEquipment.DeleteAreaRQ) (*protoEquipment.SuccessRP, error) {
	return S.c.ControllerDeleteArea(in.Id)
}

func (S *service) EditArea(ctx context.Context, in *protoEquipment.EditAreaRQ) (*protoEquipment.AreaRP, error) {
	return nil, nil
}

func (S *service) Equipment(ctx context.Context, in *protoEquipment.EquipmentRQ) (*protoEquipment.ListEquimentRP, error) {
	return S.c.ControllerEquipment(in.AreaId, in.HomeId)
}

func (S *service) CreateEquipment(ctx context.Context, in *protoEquipment.CreateEquipmentRQ) (*protoEquipment.EquipmentRP, error) {
	return S.c.ControllerCreateEquiment(in.CategoryId, in.HomeId, in.AreaId, in.Title, in.Description, in.TimeStart, in.TimeEnd, in.Cycle, in.Status)
}

func (S *service) DeleteEquipment(ctx context.Context, in *protoEquipment.DeleteEquipmentRQ) (*protoEquipment.SuccessRP, error) {
	return S.c.ControllerDeleteEquipment(in.Id)
}

func (S *service) EditEquipment(ctx context.Context, in *protoEquipment.EditEquipmentRQ) (*protoEquipment.EquipmentRP, error) {
	return nil, nil
}

func (S *service) CheckHome(ctx context.Context, in *protoEquipment.CheckHomeRQ) (*protoEquipment.HomeRP, error) {
	return S.c.ControllerCheckHome(in.AccountId, in.Id)
}
func (S *service) CheckArea(ctx context.Context, in *protoEquipment.CheckAreaRQ) (*protoEquipment.AreaRP, error) {
	return S.c.ControllercheckArea(in.Id)
}
func (S *service) CheckEquipment(ctx context.Context, in *protoEquipment.CheckEquipmentRQ) (*protoEquipment.EquipmentRP, error) {
	return S.c.ControllerCheckEquipment(in.Id)
}

func (S *service) ChangeTurnOn(ctx context.Context, in *protoEquipment.ChangeEquipmentRQ) (*protoEquipment.SuccessRP, error) {
	return S.c.ControllerChangeTurnOn(in.Id, in.TurnOn)
}

func (S *service) ChangeTime(ctx context.Context, in *protoEquipment.ChangeEquipmentTime) (*protoEquipment.SuccessRP, error) {
	return S.c.ControllerChangeTime(in.Id, in.TimeStart, in.TimeEnd)
}

func main() {
	godotenv.Load(".service.equipment.env")
	db := database.InitDB()
	// Ensure proper cleanup
	ctrl := controller.NewController(db)

	lis, err := net.Listen(os.Getenv("NET_WORK"), os.Getenv("ADDRESS")) // tao port
	if err != nil {
		log.Fatalf("err while create listen %v", err)
	}
	fmt.Println("service run on ", os.Getenv("NET_WORK"), os.Getenv("ADDRESS"))
	s := grpc.NewServer() // tao server

	protoEquipment.RegisterEquipmentServiceServer(s, &service{c: ctrl}) // dang ky server

	err = s.Serve(lis) // run server                                     // run server
	if err != nil {
		log.Fatalf("err while server %v", err)
	}
}
