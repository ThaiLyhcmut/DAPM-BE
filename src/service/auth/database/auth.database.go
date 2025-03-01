package database

import (
	"fmt"

	"github.com/ThaiLyhcmut/service/auth/model"
)

func (D *Database) CreateAccount(FullName, Email, Password, Phone string) (interface{}, error) {
	account := &model.Account{
		FullName: FullName,
		Email:    Email,
		Password: Password,
		Phone:    Phone,
	}
	result := D.DB.Create(account)
	if result.Error != nil {
		return nil, fmt.Errorf("success createAccount")
	}
	return account, nil
}
