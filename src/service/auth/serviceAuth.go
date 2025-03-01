package main

import (
	"context"
	"log"
	"net"

	protoAuth "github.com/ThaiLyhcmut/proto/auth"
	"github.com/ThaiLyhcmut/service/auth/database"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type service struct {
	protoAuth.UnimplementedAuthServiceServer
}

func (*service) Register(ctx context.Context, in *protoAuth.RegisterRQ) (*protoAuth.AccountRP, error) {
	return nil, nil
}

func (*service) Login(ctx context.Context, in *protoAuth.LoginRQ) (*protoAuth.AccountRP, error) {
	return nil, nil
}

func (*service) Infor(context.Context, *protoAuth.TokenRQ) (*protoAuth.AccountRP, error) {
	return nil, nil
}

func main() {
	godotenv.Load()
	database.InitDB()
	// Ensure proper cleanup
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

	protoAuth.RegisterAuthServiceServer(s, &service{}) // dang ky server

	err = s.Serve(lis) // run server                                     // run server
	if err != nil {
		log.Fatalf("err while server %v", err)
	}
}
