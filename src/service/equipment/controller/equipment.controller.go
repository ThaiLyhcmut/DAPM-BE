package controller

import (
	protoEquipment "ThaiLy/proto/equipment"
	"ThaiLy/service/equipment/database"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Controller struct {
	d *database.Database
}

func NewController(db *database.Database) *Controller {
	return &Controller{d: db}
}

func (C *Controller) ControllerHome(accountId int32) (*protoEquipment.ListHomeRP, error) {
	homes, err := C.d.GetHomes(accountId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "accountId invalid")
	}
	var homeRP []*protoEquipment.HomeRP
	for _, home := range homes {
		homeRP = append(homeRP, &protoEquipment.HomeRP{
			Id:       home.Id,
			HomeName: home.HomeName,
			Location: home.Location,
			Deleted:  home.Deleted,
			CreateAt: home.CreateAt.Format("2006-01-02 15:04:05"),
		})
	}
	return &protoEquipment.ListHomeRP{
		Homes: homeRP,
	}, nil
}

func (C *Controller) ControllerArea(homeId int32) (*protoEquipment.ListAreaRP, error) {
	areas, err := C.d.GetAreas(homeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "homeId invalid")
	}
	var areaRP []*protoEquipment.AreaRP
	for _, area := range areas {
		areaRP = append(areaRP, &protoEquipment.AreaRP{
			Id:     area.Id,
			HomeId: area.HomeId,
			Name:   area.Name,
		})
	}
	return &protoEquipment.ListAreaRP{
		Areas: areaRP,
	}, nil
}

func (C *Controller) ControllerEquipment(homeId int32) (*protoEquipment.ListEquimentRP, error) {
	equipments, err := C.d.GetEquipment(homeId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "homeId invalid")
	}
	var equipmentRP []*protoEquipment.EquipmentRP
	for _, equipment := range equipments {
		equipmentRP = append(equipmentRP, &protoEquipment.EquipmentRP{
			Id:          equipment.Id,
			CategoryId:  equipment.CategoryId,
			HomeId:      equipment.HomeId,
			Title:       equipment.Title,
			Description: equipment.Description,
			TimeStart:   equipment.TimeStart,
			TimeEnd:     equipment.TimeEnd,
			TurnOn:      equipment.TurnOn,
			Cycle:       equipment.Cycle,
			Status:      equipment.Status,
		})
	}
	return &protoEquipment.ListEquimentRP{
		Equipments: equipmentRP,
	}, nil
}
