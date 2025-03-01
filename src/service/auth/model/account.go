package model

type Account struct {
	Id       int    `gorm:"primaryKey"`
	FullName string `gorm:"column:FullName"`
	Email    string `gorm:"column:email"`
	Password string `gorm:"column:password"`
	Phone    string `gorm:"column:phone"`
}

func (Account) TableName() string {
	return "accounts"
}
