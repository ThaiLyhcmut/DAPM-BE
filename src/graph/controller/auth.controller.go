package controller

import (
	"ThaiLy/graph/helper"
	"ThaiLy/graph/model"
	"ThaiLy/server/client"
)

type Controller struct {
	auth *client.GRPCClient
}

// Constructor function
func NewController(auth *client.GRPCClient) *Controller {
	return &Controller{auth: auth}
}

func (C *Controller) ControllerRegister(account model.RegisterAccount) (*model.Account, error) {
	result, err := C.auth.Register(account.FullName, account.Email, account.Password, account.Phone, account.Otp)
	if err != nil {
		return nil, err
	}
	id := int(result.Id)
	token := helper.CreateJWT(id)
	resp := &model.Account{
		ID:       &id,
		FullName: &result.FullName,
		Email:    &result.Email,
		Phone:    &result.Phone,
		Token:    &token,
	}
	return resp, nil
}

func (C *Controller) ControllerLogin(account model.LoginAccount) (*model.Account, error) {
	result, err := C.auth.Login(account.Email, account.Password)
	if err != nil {
		return nil, err
	}
	id := int(result.Id)
	token := helper.CreateJWT(id)
	resp := &model.Account{
		ID:       &id,
		FullName: &result.FullName,
		Email:    &result.Email,
		Phone:    &result.Phone,
		Token:    &token,
	}
	return resp, nil
}

func (C *Controller) ControllerInfor(account model.TokenAccount) (*model.Account, error) {
	Claims, err := helper.ParseJWT(account.Token)
	if err != nil {
		return nil, err
	}
	id := int32(Claims.ID)
	result, err := C.auth.Infor(id)
	if err != nil {
		return nil, err
	}
	resp := &model.Account{
		ID:       &Claims.ID,
		FullName: &result.FullName,
		Email:    &result.Email,
		Phone:    &result.Phone,
	}
	return resp, nil
}
