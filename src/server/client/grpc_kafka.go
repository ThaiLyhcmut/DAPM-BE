package client

import (
	protoKafka "ThaiLy/proto/kafka"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCKafkaClient struct {
	conn   *grpc.ClientConn
	client protoKafka.DeviceServiceClient
}

func NewGRPCKafkaClient(addr string) (*GRPCKafkaClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	client := protoKafka.NewDeviceServiceClient(conn)
	return &GRPCKafkaClient{conn: conn, client: client}, nil
}

func (D *GRPCKafkaClient) DeviceService(ctx context.Context, req *protoKafka.DeviceRequest) (*protoKafka.DeviceResponse, error) {
	return D.client.ToggleDevice(ctx, req)
}
