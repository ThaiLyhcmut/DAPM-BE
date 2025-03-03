package database

import (
	"ThaiLy/service/equipment/model"
	"fmt"
)

func (D *Database) GetHomes(accountId int32) ([]*model.Home, error) {
	var Homes []*model.Home
	result := D.DB.Where("accountId = ?", accountId).Find(&Homes)
	if result.Error != nil {
		return nil, fmt.Errorf("error get homes")
	}
	return Homes, nil
}

func (D *Database) GetAreas(homeId int32) ([]*model.Area, error) {
	var Areas []*model.Area
	result := D.DB.Where("homeId = ?", homeId).Find(&Areas)
	if result.Error != nil {
		return nil, fmt.Errorf("error get areas")
	}
	return Areas, nil
}

func (D *Database) GetEquipment(homeId int32) ([]*model.Equipment, error) {
	var Equipments []*model.Equipment
	result := D.DB.Where("homeId = ?", homeId).Find(&Equipments)
	if result.Error != nil {
		return nil, fmt.Errorf("error get equipment")
	}
	return Equipments, nil
}
