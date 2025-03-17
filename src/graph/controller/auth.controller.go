package controller

import (
	"ThaiLy/graph/helper"
	"ThaiLy/graph/model"
	"context"
	"fmt"
)

func (C *Controller) ControllerRegister(account model.RegisterAccount) (*model.Account, error) {
	result, err := C.auth.Register(account.FullName, account.Email, account.Password, account.Phone, account.Otp)
	if err != nil {
		return nil, err
	}
	Id, err := helper.CreateAES(string(result.Id))
	if err != nil {
		return nil, err
	}
	token := helper.CreateJWT(Id)
	resp := &model.Account{
		ID:       &Id,
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
	str := fmt.Sprintf("%d", result.Id)
	ASEID, err := helper.CreateAES(str)
	if err != nil {
		return nil, err
	}
	token := helper.CreateJWT(ASEID)
	resp := &model.Account{
		ID:       &ASEID,
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
	Id, err := helper.ParseASE(Claims.ID)
	if err != nil {
		return nil, fmt.Errorf("error parse id")
	}
	result, err := C.auth.Infor(Id)
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
