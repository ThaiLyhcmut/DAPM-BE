package controller

import (
	"ThaiLy/graph/helper"
	"ThaiLy/graph/model"
	"ThaiLy/server/client"
	"context"
	"fmt"
)

// Constructor function
func NewController(auth *client.GRPCAuthClient, equipment *client.GRPCEquipmentClient) *Controller {
	return &Controller{auth: auth, equipment: equipment}
}

func (C *Controller) ControllerRegister(account model.RegisterAccount) (*model.Account, error) {
	result, err := C.auth.Register(account.FullName, account.Email, account.Password, account.Phone, account.Otp)
	if err != nil {
		return nil, err
	}
	token := helper.CreateJWT(result.Id)
	resp := &model.Account{
		ID:       &result.Id,
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
	token := helper.CreateJWT(result.Id)
	resp := &model.Account{
		ID:       &result.Id,
		FullName: &result.FullName,
		Email:    &result.Email,
		Phone:    &result.Phone,
		Token:    &token,
	}
	return resp, nil
}

func (C *Controller) ControllerInfor(ctx context.Context) (*model.Account, error) {
	Claims, ok := ctx.Value(helper.Auth).(*helper.Claims)
	if !ok {
		return nil, fmt.Errorf("Unauthorzation")
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
