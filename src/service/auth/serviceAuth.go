package main

import (
	"context"
	"io"
	"log"
	"net"

	protoAuth "ThaiLy/proto/auth"
	"ThaiLy/service/auth/controller"
	"ThaiLy/service/auth/database"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	protoAuth.UnimplementedAuthServiceServer
	c *controller.Controller
}

func (S *service) Register(stream protoAuth.AuthService_RegisterServer) error {
	var lastRequest *protoAuth.RegisterRQ

	for {
		in, err := stream.Recv()
		// đống stream
		if err == io.EOF {
			log.Println("EOF...")
			// 🔥 Nếu không có request nào gửi đến, trả lỗi
			if lastRequest == nil {
				return status.Error(codes.InvalidArgument, "No data received")
			}
			// 🔥 Gọi Controller để xử lý request cuối cùng
			resp, err := S.c.ControllerRegister(lastRequest)
			if err != nil {
				return err
			}
			// 🔥 Gửi phản hồi và đóng stream
			return stream.SendAndClose(resp)
		}
		if err != nil {
			return status.Error(codes.Internal, "Failed to receive data")
		}
		// 🔥 Cập nhật request cuối cùng nhận được
		lastRequest = in
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
