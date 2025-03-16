package controller

import (
	"ThaiLy/graph/helper"
	"ThaiLy/graph/model"
	"context"
	"fmt"
	"strconv"
)

func (C *Controller) GetHome(ctx context.Context) ([]*model.HomeQuery, error) {
	Claims, ok := ctx.Value(helper.Auth).(*helper.Claims)
	if !ok {
		return nil, fmt.Errorf("Unauthorzition")
	}
	Id, err := strconv.Atoi(helper.ParseASE(Claims.ID))
	if err != nil {
		return nil, fmt.Errorf("error parse id")
	}
	homes, err := C.equipment.GetHome(int32(Id))
	if err != nil {
		return nil, fmt.Errorf("error get home")
	}
	var resp []*model.HomeQuery
	for _, home := range homes.Homes {
		resp = append(resp, &model.HomeQuery{
			ID:        home.Id,
			HomeName:  &home.HomeName,
			Location:  &home.Location,
			Deleted:   &home.Deleted,
			CreatedAt: &home.CreatedAt,
		})
	}
	return resp, nil
}

func (C *Controller) CreateHome(ctx context.Context, home model.CreateHome) (*model.Home, error) {
	Claims, ok := ctx.Value(helper.Auth).(*helper.Claims)
	fmt.Print(Claims.ID, ok)
	if !ok {
		return nil, fmt.Errorf("Unauthorzition")
	}
	Id, err := strconv.Atoi(helper.ParseASE(Claims.ID))
	if err != nil {
		return nil, fmt.Errorf("error parse id")
	}
	Id32 := int32(Id)
	resp, err := C.equipment.CreateHome(Id32, home.HomeName, home.Location)
	if err != nil {
		return nil, fmt.Errorf("error create home")
	}
	return &model.Home{
		ID:        resp.Id,
		AccountID: &Id32,
		HomeName:  &resp.HomeName,
		Location:  &resp.Location,
		Deleted:   &resp.Deleted,
		CreatedAt: &resp.CreatedAt,
	}, nil
}

func (C *Controller) DeleteHome(ctx context.Context, home model.DeleteHome) (*model.Response, error) {
	Claims, ok := ctx.Value(helper.Auth).(*helper.Claims)
	fmt.Print(Claims.ID, ok)
	if !ok {
		return nil, fmt.Errorf("Unauthorzition")
	}
	Id, err := strconv.Atoi(helper.ParseASE(Claims.ID))
	if err != nil {
		return nil, fmt.Errorf("error parse id")
	}
	Id32 := int32(Id)
	exitsHome, err := C.equipment.CheckHome(Id32, home.ID)
	if err != nil || exitsHome == nil {
		return nil, fmt.Errorf("error check home")
	}
	resp, err := C.equipment.DeleteHome(home.ID)
	if err != nil {
		return nil, err
	}
	return &model.Response{
		Code: &resp.Code,
		Msg:  &resp.Msg,
	}, nil
}

func (C *Controller) GetArea(obj *model.HomeQuery) ([]*model.AreaQuery, error) {
	if obj == nil {
		return nil, nil
	}
	areas, err := C.equipment.GetArea(obj.ID)
	if err != nil {
		return nil, fmt.Errorf("error getArea")
	}
	var resp []*model.AreaQuery
	for _, area := range areas.Areas {
		resp = append(resp, &model.AreaQuery{
			ID:     &area.Id,
			HomeID: &area.HomeId,
			Name:   &area.Name,
		})
	}
	return resp, nil
}

func (C *Controller) CreateArea(ctx context.Context, area model.CreateArea) (*model.Area, error) {
	Claims, ok := ctx.Value(helper.Auth).(*helper.Claims)
	fmt.Print(Claims.ID, ok)
	if !ok {
		return nil, fmt.Errorf("Unauthorzition")
	}
	Id, err := strconv.Atoi(helper.ParseASE(Claims.ID))
	if err != nil {
		return nil, fmt.Errorf("error parse id")
	}
	Id32 := int32(Id)
	exitsHome, err := C.equipment.CheckHome(Id32, area.HomeID)
	if err != nil || exitsHome == nil {
		return nil, fmt.Errorf("not exits home")
	}
	resp, err := C.equipment.CreateArea(area.HomeID, area.Name)
	if err != nil {
		return nil, fmt.Errorf("error create area")
	}
	return &model.Area{
		ID:     &resp.Id,
		HomeID: &resp.HomeId,
		Name:   &resp.Name,
	}, nil
}

func (C *Controller) DeleteArea(ctx context.Context, area model.DeleteArea) (*model.Response, error) {
	Claims, ok := ctx.Value(helper.Auth).(*helper.Claims)
	fmt.Print(Claims.ID, ok)
	if !ok {
		return nil, fmt.Errorf("Unauthorzition")
	}
	exitsArea, err := C.equipment.CheckArea(area.ID)
	if err != nil || exitsArea == nil {
		return nil, fmt.Errorf("error check Area")
	}
	Id, err := strconv.Atoi(helper.ParseASE(Claims.ID))
	if err != nil {
		return nil, fmt.Errorf("error parse id")
	}
	Id32 := int32(Id)
	exitsHome, err := C.equipment.CheckHome(Id32, exitsArea.HomeId)
	if err != nil || exitsHome == nil {
		return nil, fmt.Errorf("error check home")
	}
	resp, err := C.equipment.DeleteArea(area.ID)
	if err != nil {
		return nil, err
	}
	return &model.Response{
		Code: &resp.Code,
		Msg:  &resp.Msg,
	}, nil
}

func (C *Controller) GetEquipment(obj *model.AreaQuery) ([]*model.Equipment, error) {
	equipments, err := C.equipment.GetEquipment(*obj.ID, *obj.HomeID)
	if err != nil {
		return nil, fmt.Errorf("error get equipment")
	}
	var resp []*model.Equipment
	for _, equipment := range equipments.Equipments {
		resp = append(resp, &model.Equipment{
			ID:          &equipment.Id,
			CategoryID:  &equipment.CategoryId,
			HomeID:      &equipment.HomeId,
			AreaID:      &equipment.AreaId,
			Title:       &equipment.Title,
			Description: &equipment.Description,
			TimeStart:   &equipment.TimeStart,
			TimeEnd:     &equipment.TimeEnd,
			TurnOn:      &equipment.TurnOn,
			Cycle:       &equipment.Cycle,
			Status:      &equipment.Status,
		})
	}
	return resp, nil
}

func (C *Controller) CreateEquiment(ctx context.Context, equipment model.CreateEquiment) (*model.Equipment, error) {
	Claims, ok := ctx.Value(helper.Auth).(*helper.Claims)
	fmt.Print(Claims.ID, ok)
	if !ok {
		return nil, fmt.Errorf("Unauthorzition")
	}
	Id, err := strconv.Atoi(helper.ParseASE(Claims.ID))
	if err != nil {
		return nil, fmt.Errorf("error parse id")
	}
	Id32 := int32(Id)
	exitsHome, err := C.equipment.CheckHome(Id32, equipment.HomeID)
	if exitsHome == nil || err != nil {
		return nil, fmt.Errorf("not exits home")
	}
	exitsArea, err := C.equipment.CheckArea(equipment.AreaID)
	if exitsArea == nil || err != nil {
		return nil, fmt.Errorf("not exits area")
	}
	if exitsArea.HomeId != equipment.HomeID {
		return nil, fmt.Errorf(":)) bạn đừng đùa vậy")
	}
	resp, err := C.equipment.CreateEquipment(equipment.CategoryID, equipment.HomeID, equipment.AreaID, equipment.Title, *equipment.Description, equipment.Status)
	if err != nil {
		return nil, fmt.Errorf("error create area")
	}
	return &model.Equipment{
		ID:          &resp.Id,
		CategoryID:  &resp.CategoryId,
		HomeID:      &resp.HomeId,
		AreaID:      &resp.AreaId,
		Title:       &resp.Title,
		Description: &resp.Description,
		TimeStart:   &resp.TimeStart,
		TimeEnd:     &resp.TimeEnd,
		TurnOn:      &resp.TurnOn,
		Cycle:       &resp.Cycle,
		Status:      &resp.Status,
	}, nil
}

func (C *Controller) DeleteEquipment(ctx context.Context, equipment model.DeleteEquipment) (*model.Response, error) {
	Claims, ok := ctx.Value(helper.Auth).(*helper.Claims)
	fmt.Print(Claims.ID, ok)
	if !ok {
		return nil, fmt.Errorf("Unauthorzition")
	}
	exitsEquipment, err := C.equipment.CheckEquipment(equipment.ID)
	if err != nil || exitsEquipment == nil {
		return nil, fmt.Errorf("error check Area")
	}
	Id, err := strconv.Atoi(helper.ParseASE(Claims.ID))
	if err != nil {
		return nil, fmt.Errorf("error parse id")
	}
	Id32 := int32(Id)
	exitsHome, err := C.equipment.CheckHome(Id32, exitsEquipment.HomeId)
	if err != nil || exitsHome == nil {
		return nil, fmt.Errorf("error check home")
	}
	resp, err := C.equipment.DeleteEquipment(equipment.ID)
	if err != nil {
		return nil, err
	}
	return &model.Response{
		Code: &resp.Code,
		Msg:  &resp.Msg,
	}, nil
}
