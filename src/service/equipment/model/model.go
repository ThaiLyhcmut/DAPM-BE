package model

type Home struct {
	Id        int32  `gorm:"primarykey"`
	AccountId int32  `gorm:"column:accountId"`
	HomeName  string `gorm:"column:homeName"`
	Location  string `gorm:"column:location"`
	Deleted   bool   `gorm:"column:deleted"`
	CreatedAt string `gorm:"column:createdAt"`
}

func (Home) TableName() string {
	return "myHome"
}

type Area struct {
	Id     int32  `gorm:"primarykey"`
	HomeId int32  `gorm:"column:homeId"`
	Name   string `gorm:"column:name"`
}

func (Area) TableName() string {
	return "areas"
}

type Equipment struct {
	Id          int32  `gorm:"primarykey"`
	CategoryId  int32  `gorm:"column:categoryId"`
	HomeId      int32  `gorm:"column:homeId"`
	AreaId      int32  `gorm:"column:areaId"`
	Title       string `gorm:"column:title"`
	Description string `gorm:"column:description"`
	TimeStart   string `gorm:"column:timeStart"`
	TimeEnd     string `gorm:"column:timeEnd"`
	TurnOn      bool   `gorm:"column:turnOn"`
	Cycle       int32  `gorm:"column:cycle"`
	Status      string `gorm:"column:status"`
}

func (Equipment) TableName() string {
	return "equipments"
}
