package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"

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
	//
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

func (S *service) Login(ctx context.Context, in *protoAuth.LoginRQ) (*protoAuth.AccountRP, error) {
	return S.c.ControllerLogin(in)
}

func (S *service) Infor(ctx context.Context, in *protoAuth.IdA) (*protoAuth.AccountRP, error) {
	return S.c.ControllerInfor(in)
}

func main() {
	godotenv.Load(".service.auth.env")
	db := database.InitDB()
	// Ensure proper cleanup
	ctrl := controller.NewController(db)

	lis, err := net.Listen(os.Getenv("NET_WORK"), os.Getenv("ADDRESS")) // tao port
	if err != nil {
		log.Fatalf("err while create listen %v", err)
	}
	fmt.Println("service run on ", os.Getenv("NET_WORK"), os.Getenv("ADDRESS"))

	s := grpc.NewServer() // tao server

	protoAuth.RegisterAuthServiceServer(s, &service{c: ctrl}) // dang ky server

	err = s.Serve(lis) // run server                                     // run server
	if err != nil {
		log.Fatalf("err while server %v", err)
	}
}
