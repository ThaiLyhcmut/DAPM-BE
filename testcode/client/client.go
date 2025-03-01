package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ThaiLyhcmut/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cc, err := grpc.NewClient("localhost:50069", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("err while dial %v", err)
	}
	defer cc.Close()
	client := calculatorpb.NewCalculatorServiceClient(cc)
	// log.Printf("service client %f", client)
	callSum(client)
}

func callSum(c calculatorpb.CalculatorServiceClient) {
	resp, err := c.Sum(context.Background(), &calculatorpb.SumRequest{
		Num1: 5,
		Num2: 10,
	})
	if err != nil {
		log.Fatal("call error %v", err)
	}
	fmt.Println(resp)
}
