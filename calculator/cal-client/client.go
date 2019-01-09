package main

import (
	"context"
	"fmt"
	"io"
	"log"

	calpb "github.com/grpc-go-course/calculator/pb"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatal("Failed to initial client")
	}

	defer cc.Close()

	calClient := calpb.NewCalculatorClient(cc)

	// log.Println(calSum(calClient))
	calPrime(calClient)
}

func calPrime(calClient calpb.CalculatorClient) {
	req := &calpb.PrimeRequest{
		PrimeNum: 120,
	}

	clienStream, err := calClient.GetPrimeNum(context.Background(), req)
	if err != nil {
		log.Fatal("Failed to get prime number")
	}
	for {
		res, err := clienStream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal("Failed get response")
		}
		fmt.Print(" ", res.GetPrimeResult())
	}
}

func calSum(calClient calpb.CalculatorClient) int32 {
	req := calpb.SumRequest{
		Num1: 10,
		Num2: 3,
	}

	res, err := calClient.GetSum(context.Background(), &req)
	if err != nil {
		log.Fatal("Failed to try calculating sum of ", req.Num1, req.Num2)
	}
	log.Printf("Sum of %v %v is %v ", req.Num1, req.Num2, res.GetSum())
	return res.Sum
}
