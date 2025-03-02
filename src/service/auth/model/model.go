package model

import "time"

type Account struct {
	Id       int32  `gorm:"primaryKey"`
	FullName string `gorm:"column:FullName"`
	Email    string `gorm:"column:email"`
	Password string `gorm:"column:password"`
	Phone    string `gorm:"column:phone"`
}

func (Account) TableName() string {
	return "accounts"
}

type Otp struct {
	Id        int       `gorm:"primaryKey"`
	Email     string    `gorm:"column:email"`
	Code      string    `gorm:"column:code"`
	ExpiredAt time.Time `gorm:"column:expiredAt"`
	Used      int       `gorm:"column:used"`
}
