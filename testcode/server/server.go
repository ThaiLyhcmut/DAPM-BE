package main

import (
	"context"
	"io"
	"log"
	"math"
	"net"
	"time"

	"github.com/ThaiLyhcmut/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, in *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	resp := &calculatorpb.SumResponse{
		Result: in.GetNum1() + in.GetNum2(),
	}
	return resp, nil
}

func (*server) PrimeNumberDecompostion(in *calculatorpb.PNDRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompostionServer) error {
	n := in.GetNumber()
	var k int32 = 2

	for n > 1 {
		if n%k == 0 {
			resp := &calculatorpb.PNDResponse{Result: k}
			if err := stream.Send(resp); err != nil {
				return err
			}
			n = n / k
			time.Sleep(5 * time.Second)
		} else {
			k++
		}
	}

	return nil
}

func (*server) Average(stream grpc.BidiStreamingServer[calculatorpb.AVARequest, calculatorpb.AVAResponse]) error {
	var total float32 = 0
	var count int = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			resp := &calculatorpb.AVAResponse{
				Result: total / float32(count),
			}
			return stream.Send(resp)
		}
		if err != nil {
			log.Fatal("err %v", err)
		}
		total += req.GetNum()
		count++
	}
	return nil
}

func (*server) Max(stream grpc.BidiStreamingServer[calculatorpb.MRequest, calculatorpb.MResponse]) error {
	var MAX int32 = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("EOF...")
			return nil
		}
		if err != nil {
			log.Fatal("err %v", err)
		}
		MAX = max(MAX, req.GetNum())
		resp := &calculatorpb.MResponse{
			Max: MAX,
		}
		err = stream.Send(resp)
		if err != nil {
			return err
		}
	}
}

func (*server) Squere(ctx context.Context, in *calculatorpb.SQRequest) (*calculatorpb.SQResponse, error) {
	r := in.GetNum()
	if r < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Expect R > 0, req R was %v", r)
	}
	return &calculatorpb.SQResponse{
		SQRoot: math.Sqrt(float64(r)),
	}, nil
}

func (*server) SumWithDeadline(ctx context.Context, in *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			return nil, status.Errorf(codes.Canceled, "Client canceled the request")
		}
		if ctx.Err() == context.DeadlineExceeded {
			return nil, status.Errorf(codes.DeadlineExceeded, "Deadline exceeded")
		}
		time.Sleep(1 * time.Second)
	}
	resp := &calculatorpb.SumResponse{
		Result: in.GetNum1() + in.GetNum2(),
	}
	return resp, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50069") // tao port
	if err != nil {
		log.Fatalf("err while create listen %v", err)
	}

	certFile := "ssl/server.crt"
	keyFile := "ssl/server.pem"
	//ssl for server
	creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
	if sslErr != nil {
		log.Fatal("create creds ssl err %v\n", sslErr)
	}

	opts := grpc.Creds(creds)

	s := grpc.NewServer(opts) // tao server

	calculatorpb.RegisterCalculatorServiceServer(s, &server{}) // dang ky server

	err = s.Serve(lis) // run server                                     // run server
	if err != nil {
		log.Fatalf("err while server %v", err)
	}
}
