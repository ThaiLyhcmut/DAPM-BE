package database

import (
	"fmt"

	"ThaiLy/service/auth/model"
)

func (D *Database) CreateAccount(FullName, Email, Password, Phone string) (*model.Account, error) {
	account := &model.Account{
		FullName: FullName,
		Email:    Email,
		Password: Password,
		Phone:    Phone,
	}
	result := D.DB.Create(account)
	if result.Error != nil {
		return nil, fmt.Errorf("error createAccount")
	}
	return account, nil
}

func (D *Database) LoginAccount(Email, Password string) (*model.Account, error) {
	account := &model.Account{}
	result := D.DB.Where("email = ? AND password = ?", Email, Password).First(account)
	if result.Error != nil {
		return nil, fmt.Errorf("error login")
	}
	return account, nil
}

func (D *Database) InforAccount(Id int) (*model.Account, error) {
	account := &model.Account{}
	result := D.DB.First(account, Id)
	if result.Error != nil {
		return nil, fmt.Errorf("error infor")
	}
	return account, nil
}

func (D *Database) DeleteOtp(Email, Otp string) error {
	result := D.DB.Where(&model.Otp{Email: Email, Code: Otp}).First(&model.Otp{})
	if result.Error != nil {
		return fmt.Errorf("error deleteOtp")
	}
	return nil
}
