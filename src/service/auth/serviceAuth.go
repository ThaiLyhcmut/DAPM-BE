package main

import (
	"context"
	"io"
	"log"
	"net"

	protoAuth "github.com/ThaiLyhcmut/proto/auth"
	"github.com/ThaiLyhcmut/service/auth/controller"
	"github.com/ThaiLyhcmut/service/auth/database"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	protoAuth.UnimplementedAuthServiceServer
	c *controller.Controller
}

func (S *service) Register(stream grpc.ClientStreamingServer[protoAuth.RegisterRQ, protoAuth.AccountRP]) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			log.Println("EOF...")
			return nil
		}
		if err != nil {
			return status.Error(codes.InvalidArgument, "input invalid")
		}
		if in.GetOtp() != "" {
			resp, err := S.c.ControllerRegister(in)
			if err != nil {
				return err
			}
			return stream.SendAndClose(resp)
		}

	}
}

func (*service) Login(ctx context.Context, in *protoAuth.LoginRQ) (*protoAuth.AccountRP, error) {
	return nil, nil
}

func (*service) Infor(context.Context, *protoAuth.TokenRQ) (*protoAuth.AccountRP, error) {
	return nil, nil
}

func main() {
	godotenv.Load()
	db := database.InitDB()
	// Ensure proper cleanup
	ctrl := controller.NewController(db)

	lis, err := net.Listen("tcp", "localhost:55555") // tao port
	if err != nil {
		log.Fatalf("err while create listen %v", err)
	}

	// certFile := "ssl/server.crt"
	// keyFile := "ssl/server.pem"
	// //ssl for server
	// creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
	// if sslErr != nil {
	// 	log.Fatal("create creds ssl err %v\n", sslErr)
	// }

	// opts := grpc.Creds(creds)

	s := grpc.NewServer() // tao server

	protoAuth.RegisterAuthServiceServer(s, &service{c: ctrl}) // dang ky server

	err = s.Serve(lis) // run server                                     // run server
	if err != nil {
		log.Fatalf("err while server %v", err)
	}
}
