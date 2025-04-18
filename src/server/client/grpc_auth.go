package client

import (
	"context"
	"fmt"

	protoAuth "ThaiLy/proto/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCAuthClient struct {
	conn   *grpc.ClientConn
	client protoAuth.AuthServiceClient
}

func NewGRPCAuthClient(addr string) (*GRPCAuthClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	client := protoAuth.NewAuthServiceClient(conn)
	return &GRPCAuthClient{conn: conn, client: client}, nil
}

func (c *GRPCAuthClient) Register(fullName string, email string, password string, phone string, otp *string) (*protoAuth.AccountRP, error) {
	// 1. Mở stream
	stream, err := c.client.Register(context.Background())
	if err != nil {
		return nil, err
	}
	var isOtp string = ""
	if otp != nil {
		isOtp = *otp
	}

	// 2. Gửi dữ liệu
	req := &protoAuth.RegisterRQ{
		FullName: fullName,
		Email:    email,
		Password: password,
		Phone:    phone,
		Otp:      isOtp,
	}
	fmt.Println(req)
	if err := stream.Send(req); err != nil {
		return nil, err
	}

	// 3. Đóng stream và nhận phản hồi từ server
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return nil, err
	}
	fmt.Println(resp)

	return resp, nil
}

func (c *GRPCAuthClient) Login(email, password string) (*protoAuth.AccountRP, error) {
	in := &protoAuth.LoginRQ{
		Email:    email,
		Password: password,
	}
	return c.client.Login(context.Background(), in)
}

func (c *GRPCAuthClient) Infor(id int32) (*protoAuth.AccountRP, error) {
	in := &protoAuth.IdA{
		Id: id,
	}
	return c.client.Infor(context.Background(), in)
}
