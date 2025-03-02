package controller

import (
	protoAuth "ThaiLy/proto/auth"
	"fmt"

	"ThaiLy/service/auth/database"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Controller struct {
	d *database.Database
}

func NewController(db *database.Database) *Controller {
	return &Controller{d: db}
}

func (C *Controller) ControllerRegister(in *protoAuth.RegisterRQ) (*protoAuth.AccountRP, error) {
	if result := C.d.DeleteOtp(in.Email, in.Otp); result != nil {
		return nil, status.Error(codes.InvalidArgument, "email and otp invalid")
	}
	account, err := C.d.CreateAccount(in.FullName, in.Email, in.Password, in.Phone)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "email password phone invalid")
	}
	return &protoAuth.AccountRP{
		Id:       account.Id,
		FullName: account.FullName,
		Email:    account.Email,
		Phone:    account.Phone,
	}, nil
}

func (C *Controller) ControllerLogin(in *protoAuth.LoginRQ) (*protoAuth.AccountRP, error) {
	account, err := C.d.LoginAccount(in.Email, in.Password)
	fmt.Println(account, err)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "email password invalid")
	}
	return &protoAuth.AccountRP{
		Id:       account.Id,
		FullName: account.FullName,
		Email:    account.Email,
		Phone:    account.Phone,
	}, nil
}

func (C *Controller) ControllerInfor(in *protoAuth.IdA) (*protoAuth.AccountRP, error) {
	account, err := C.d.InforAccount(int(in.GetId()))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "id invalid")
	}
	return &protoAuth.AccountRP{
		Id:       account.Id,
		FullName: account.FullName,
		Email:    account.Email,
		Phone:    account.Phone,
	}, nil
}
