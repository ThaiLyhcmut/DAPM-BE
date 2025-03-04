package controller

import (
	protoEquipment "ThaiLy/proto/equipment"
	"ThaiLy/service/equipment/database"
	"ThaiLy/service/equipment/model"

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
			Id:        home.Id,
			HomeName:  home.HomeName,
			Location:  home.Location,
			Deleted:   home.Deleted,
			CreatedAt: home.CreatedAt,
		})
	}
	return &protoEquipment.ListHomeRP{
		Homes: homeRP,
	}, nil
}

func (C *Controller) ControllerCreateHome(accountId int32, homeName string, location string, deleted bool, createdAt string) (*protoEquipment.HomeRP, error) {
	home, err := C.d.CreateHome(accountId, homeName, location, deleted, createdAt)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "create home data invalid")
	}
	return &protoEquipment.HomeRP{
		Id:        home.Id,
		HomeName:  home.HomeName,
		Location:  home.Location,
		Deleted:   home.Deleted,
		CreatedAt: home.CreatedAt,
	}, nil
}

func (C *Controller) ControllerDeleteHome(id int32) (*protoEquipment.SuccessRP, error) {
	err := C.d.DeleteHome(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "delete home invalid")
	}
	return &protoEquipment.SuccessRP{
		Code: "200",
		Msg:  "success",
	}, nil
}

func (C *Controller) ControllerCheckHome(accountId int32, id int32) (*protoEquipment.HomeRP, error) {
	home, err := C.d.CheckHome(accountId, id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "check home invalid")
	}
	return &protoEquipment.HomeRP{
		Id:        home.Id,
		HomeName:  home.HomeName,
		Location:  home.Location,
		Deleted:   home.Deleted,
		CreatedAt: home.CreatedAt,
	}, nil
}
func (C *Controller) ControllercheckArea(id int32) (*protoEquipment.AreaRP, error) {
	area, err := C.d.CheckArea(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "check home invalid")
	}
	return &protoEquipment.AreaRP{
		Id:     area.Id,
		HomeId: area.HomeId,
		Name:   area.Name,
	}, nil
}
func (C *Controller) ControllerCheckEquipment(id int32) (*protoEquipment.EquipmentRP, error) {
	equipment, err := C.d.CheckEquipment(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "check home invalid")
	}
	return &protoEquipment.EquipmentRP{
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

func (C *Controller) ControllerCreateArea(homeId int32, name string) (*protoEquipment.AreaRP, error) {
	area, err := C.d.CreateArea(homeId, name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "error create area")
	}
	return &protoEquipment.AreaRP{
		Id:     area.Id,
		HomeId: area.HomeId,
		Name:   area.Name,
	}, nil
}

func (C *Controller) ControllerDeleteArea(id int32) (*protoEquipment.SuccessRP, error) {
	err := C.d.DeleteArea(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "error delete area")
	}
	return &protoEquipment.SuccessRP{
		Code: "200",
		Msg:  "success",
	}, nil
}

func (C *Controller) ControllerEquipment(areaId int32, homeId int32) (*protoEquipment.ListEquimentRP, error) {
	var equipments []*model.Equipment
	var err error
	if areaId == 0 {
		equipments, err = C.d.GetEquipmentByHomeId(homeId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "homeId invalid")
		}
	} else {
		equipments, err = C.d.GetEquipment(areaId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "homeId invalid")
		}
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

func (C *Controller) ControllerCreateEquiment(categoryId int32, homeId int32, title string, description string, timeStart string, timeEnd string, cycle int32, stats string) (*protoEquipment.EquipmentRP, error) {
	equipment, err := C.d.CreateEquipment(categoryId, homeId, title, description, timeStart, timeEnd, cycle, stats)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "homeId invalid")
	}
	return &protoEquipment.EquipmentRP{
		Id:          equipment.Id,
		CategoryId:  equipment.CategoryId,
		HomeId:      equipment.HomeId,
		AreaId:      equipment.AreaId,
		Title:       equipment.Title,
		Description: equipment.Description,
		TimeStart:   equipment.TimeStart,
		TimeEnd:     equipment.TimeEnd,
		Cycle:       equipment.Cycle,
		Status:      equipment.Status,
	}, nil
}

func (C *Controller) ControllerDeleteEquipment(id int32) (*protoEquipment.SuccessRP, error) {
	err := C.d.DeleteEquipment(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "error delete equipment")
	}
	return &protoEquipment.SuccessRP{
		Code: "200",
		Msg:  "success",
	}, nil
}
