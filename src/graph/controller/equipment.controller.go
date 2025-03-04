package controller

import (
	"ThaiLy/graph/helper"
	"ThaiLy/graph/model"
	"context"
	"fmt"
)

func (C *Controller) GetHome(ctx context.Context) ([]*model.Home, error) {
	Claims, ok := ctx.Value(helper.Auth).(*helper.Claims)
	fmt.Print(Claims.ID, ok)
	if !ok {
		return nil, fmt.Errorf("Unauthorzition")
	}
	homes, err := C.equipment.GetHome(Claims.ID)
	if err != nil {
		return nil, fmt.Errorf("error get home")
	}
	var resp []*model.Home
	for _, home := range homes.Homes {
		resp = append(resp, &model.Home{
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
	resp, err := C.equipment.CreateHome(Claims.ID, home.HomeName, home.Location)
	if err != nil {
		return nil, fmt.Errorf("error create home")
	}
	return &model.Home{
		ID:        resp.Id,
		AccountID: &Claims.ID,
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
	exitsHome, err := C.equipment.CheckHome(Claims.ID, home.ID)
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

func (C *Controller) GetArea(obj *model.Home) ([]*model.Area, error) {
	if obj == nil {
		return nil, nil
	}
	areas, err := C.equipment.GetArea(obj.ID)
	if err != nil {
		return nil, fmt.Errorf("error getArea")
	}
	var resp []*model.Area
	for _, area := range areas.Areas {
		resp = append(resp, &model.Area{
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
	// check nha
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
	exitsHome, err := C.equipment.CheckHome(Claims.ID, exitsArea.HomeId)
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

func (C *Controller) GetEquipment(obj *model.Area) ([]*model.Equipment, error) {
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
	// check nha
	// check khu vuc
	resp, err := C.equipment.CreateEquipment(equipment.CategoryID, equipment.HomeID, equipment.Title, *equipment.Description, equipment.Status)
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
	exitsHome, err := C.equipment.CheckHome(Claims.ID, exitsEquipment.HomeId)
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
