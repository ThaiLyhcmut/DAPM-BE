package client

import (
	protoEquipment "ThaiLy/proto/equipment"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCEquipmentClient struct {
	conn   *grpc.ClientConn
	client protoEquipment.EquipmentServiceClient
}

func NewGRPCEquipmentClient(addr string) (*GRPCEquipmentClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	client := protoEquipment.NewEquipmentServiceClient(conn)
	return &GRPCEquipmentClient{conn: conn, client: client}, nil
}

func (c *GRPCEquipmentClient) GetHome(accountId int32) (*protoEquipment.ListHomeRP, error) {
	return c.client.Home(context.Background(), &protoEquipment.HomeRQ{
		AccountId: accountId,
	})
}

func (c *GRPCEquipmentClient) CreateHome(accountId int32, homeName string, location string) (*protoEquipment.HomeRP, error) {
	return c.client.CreateHome(context.Background(), &protoEquipment.CreateHomeRQ{
		AccountId: accountId,
		HomeName:  homeName,
		Location:  location,
	})
}

func (c *GRPCEquipmentClient) CheckHome(accountId int32, id int32) (*protoEquipment.HomeRP, error) {
	return c.client.CheckHome(context.Background(), &protoEquipment.CheckHomeRQ{
		AccountId: accountId,
		Id:        id,
	})
}
func (c *GRPCEquipmentClient) CheckArea(id int32) (*protoEquipment.AreaRP, error) {
	return c.client.CheckArea(context.Background(), &protoEquipment.CheckAreaRQ{
		Id: id,
	})
}
func (c *GRPCEquipmentClient) CheckEquipment(id int32) (*protoEquipment.EquipmentRP, error) {
	return c.client.CheckEquipment(context.Background(), &protoEquipment.CheckEquipmentRQ{
		Id: id,
	})
}

func (c *GRPCEquipmentClient) DeleteHome(id int32) (*protoEquipment.SuccessRP, error) {
	return c.client.DeleteHome(context.Background(), &protoEquipment.DeleteHomeRQ{
		Id: id,
	})
}

func (c *GRPCEquipmentClient) GetArea(homeId int32) (*protoEquipment.ListAreaRP, error) {
	return c.client.Area(context.Background(), &protoEquipment.AreaRQ{
		HomeId: homeId,
	})
}

func (c *GRPCEquipmentClient) CreateArea(homeId int32, name string) (*protoEquipment.AreaRP, error) {
	return c.client.CreateArea(context.Background(), &protoEquipment.CreateAreaRQ{
		HomeId: homeId,
		Name:   name,
	})
}

func (c *GRPCEquipmentClient) DeleteArea(id int32) (*protoEquipment.SuccessRP, error) {
	return c.client.DeleteArea(context.Background(), &protoEquipment.DeleteAreaRQ{
		Id: id,
	})
}

func (c *GRPCEquipmentClient) GetEquipment(areaId int32, homeId int32) (*protoEquipment.ListEquimentRP, error) {
	return c.client.Equipment(context.Background(), &protoEquipment.EquipmentRQ{
		AreaId: areaId,
		HomeId: homeId,
	})
}

func (c *GRPCEquipmentClient) CreateEquipment(categoryId int32, homeId int32, areaId int32, title string, description string, status string) (*protoEquipment.EquipmentRP, error) {
	return c.client.CreateEquipment(context.Background(), &protoEquipment.CreateEquipmentRQ{
		CategoryId:  categoryId,
		HomeId:      homeId,
		AreaId:      areaId,
		Title:       title,
		Description: description,
		Status:      status,
	})
}

func (c *GRPCEquipmentClient) DeleteEquipment(id int32) (*protoEquipment.SuccessRP, error) {
	return c.client.DeleteEquipment(context.Background(), &protoEquipment.DeleteEquipmentRQ{
		Id: id,
	})
}
