package database

import (
	"ThaiLy/service/equipment/model"
	"fmt"
	"time"
)

func (D *Database) GetHomes(accountId int32) ([]*model.Home, error) {
	var Homes []*model.Home
	result := D.DB.Where("accountId = ?", accountId).Find(&Homes)
	if result.Error != nil {
		return nil, fmt.Errorf("error get homes")
	}
	return Homes, nil
}

func (D *Database) CreateHome(accountId int32, homeName string, location string, deleted bool, createdAt string) (*model.Home, error) {
	Home := &model.Home{
		AccountId: accountId,
		HomeName:  homeName,
		Location:  location,
		Deleted:   deleted,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"), // Định dạng chuẩn của Go
	}
	fmt.Println(*Home)
	result := D.DB.Create(Home)
	if result.Error != nil {
		return nil, fmt.Errorf("error create home")
	}
	return Home, nil
}

func (D *Database) DeleteHome(id int32) error {
	result := D.DB.Delete(&model.Home{}, id)
	if result.Error != nil {
		return fmt.Errorf("error delete home")
	}
	return nil
}

func (D *Database) EditHome(id int32, homeName string, location string, deleted bool) (*model.Home, error) {
	home := &model.Home{
		Id: id,
	}
	result := D.DB.First(home)
	if result.Error != nil {
		return nil, fmt.Errorf("err edit home")
	}
	home.HomeName = homeName
	home.Location = location
	home.Deleted = deleted
	D.DB.Save(home)
	return home, nil
}

func (D *Database) CheckHome(accountId int32, id int32) (*model.Home, error) {
	Home := &model.Home{}
	result := D.DB.Where("accountId = ? AND id = ?", accountId, id).First(Home)
	if result.Error != nil {
		return nil, fmt.Errorf("error check home")
	}
	return Home, nil
}

func (D *Database) GetAreas(homeId int32) ([]*model.Area, error) {
	var Areas []*model.Area
	result := D.DB.Where("homeId = ?", homeId).Find(&Areas)
	if result.Error != nil {
		return nil, fmt.Errorf("error get areas")
	}
	return Areas, nil
}

func (D *Database) CreateArea(homeId int32, name string) (*model.Area, error) {
	Area := &model.Area{
		HomeId: homeId,
		Name:   name,
	}
	result := D.DB.Create(Area)
	if result.Error != nil {
		return nil, fmt.Errorf("error create area")
	}
	return Area, nil
}

func (D *Database) DeleteArea(id int32) error {
	result := D.DB.Delete(&model.Area{}, id)
	if result.Error != nil {
		return fmt.Errorf("error delete area")
	}
	return nil
}

func (D *Database) EditArea(id int32, homeId int32, name string) (*model.Area, error) {
	area := &model.Area{
		Id: id,
	}
	result := D.DB.First(area)
	if result.Error != nil {
		return nil, fmt.Errorf("err edit area")
	}
	area.HomeId = homeId
	area.Name = name
	D.DB.Save(area)
	return area, nil
}

func (D *Database) CheckArea(id int32) (*model.Area, error) {
	area := &model.Area{
		Id: id,
	}
	result := D.DB.First(area)
	if result.Error != nil {
		return nil, fmt.Errorf("error check area")
	}
	return area, nil
}

func (D *Database) GetEquipment(areaId int32) ([]*model.Equipment, error) {
	var Equipments []*model.Equipment
	result := D.DB.Where("areaId = ?", areaId).Find(&Equipments)
	if result.Error != nil {
		return nil, fmt.Errorf("error get equipment")
	}
	return Equipments, nil
}

func (D *Database) GetEquipmentByHomeId(homeId int32) ([]*model.Equipment, error) {
	var Equipments []*model.Equipment
	result := D.DB.Where("homeId = ?", homeId).Find(&Equipments)
	if result.Error != nil {
		return nil, fmt.Errorf("error get equipment")
	}
	return Equipments, nil
}

func (D *Database) CreateEquipment(categoryId int32, homeId int32, areaId int32, title string, description string, timeStart string, timeEnd string, cycle int32, status string) (*model.Equipment, error) {
	Equiment := &model.Equipment{
		CategoryId:  categoryId,
		HomeId:      homeId,
		AreaId:      areaId,
		Title:       title,
		Description: description,
		TimeStart:   time.Now().Format("2006-01-02 15:04:05"),
		TimeEnd:     time.Now().Format("2006-01-02 15:04:05"),
		Cycle:       0,
		Status:      status,
	}
	result := D.DB.Create(Equiment)
	if result.Error != nil {
		return nil, fmt.Errorf("error create equipment")
	}
	return Equiment, nil
}

func (D *Database) DeleteEquipment(id int32) error {
	result := D.DB.Delete(&model.Equipment{}, id)
	if result.Error != nil {
		return fmt.Errorf("error delete equipment")
	}
	return nil
}

func (D *Database) CheckEquipment(id int32) (*model.Equipment, error) {
	equipment := &model.Equipment{
		Id: id,
	}
	result := D.DB.First(equipment)
	if result.Error != nil {
		return nil, fmt.Errorf("error check equipment")
	}
	return equipment, nil
}
