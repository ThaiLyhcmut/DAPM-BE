package controller

import (
	"ThaiLy/graph/helper"
	"ThaiLy/graph/model"
	"context"
	"fmt"
	"strconv"
)

func (C *Controller) ControllerRegister(account model.RegisterAccount) (*model.Account, error) {
	result, err := C.auth.Register(account.FullName, account.Email, account.Password, account.Phone, account.Otp)
	if err != nil {
		return nil, err
	}
	Id := helper.CreateASE(string(result.Id))
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
	ASEID := helper.CreateASE(string(result.Id))
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
	id, err := strconv.Atoi(helper.ParseASE(Claims.ID))
	if err != nil {
		return nil, fmt.Errorf("error parse id")
	}
	result, err := C.auth.Infor(int32(id))
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
