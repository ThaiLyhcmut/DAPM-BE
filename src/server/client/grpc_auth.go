package client

import (
	"context"

	protoAuth "ThaiLy/proto/auth"

	"google.golang.org/grpc"
)

type GRPCClient struct {
	conn   *grpc.ClientConn
	client protoAuth.AuthServiceClient
}

func NewGRPCClient(addr string) (*GRPCClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := protoAuth.NewAuthServiceClient(conn)
	return &GRPCClient{conn: conn, client: client}, nil
}

func (c *GRPCClient) Register(fullName, email, password, phone string) (*protoAuth.AccountRP, error) {
	// 1. Mở stream
	stream, err := c.client.Register(context.Background())
	if err != nil {
		return nil, err
	}

	// 2. Gửi dữ liệu
	req := &protoAuth.RegisterRQ{
		FullName: fullName,
		Email:    email,
		Password: password,
		Phone:    phone,
	}

	if err := stream.Send(req); err != nil {
		return nil, err
	}

	// 3. Đóng stream và nhận phản hồi từ server
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return nil, err
	}

	return resp, nil
}
