package client

import (
	"context"
	"fmt"

	protoAuth "ThaiLy/proto/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	conn   *grpc.ClientConn
	client protoAuth.AuthServiceClient
}

func NewGRPCClient(addr string) (*GRPCClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	client := protoAuth.NewAuthServiceClient(conn)
	return &GRPCClient{conn: conn, client: client}, nil
}

func (c *GRPCClient) Register(fullName string, email string, password string, phone string, otp *string) (*protoAuth.AccountRP, error) {
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
