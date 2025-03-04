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
