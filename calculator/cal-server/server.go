package main

import (
	"context"
	"io"
	"log"
	"net"

	calpb "github.com/grpc-go-course/calculator/pb"
	"google.golang.org/grpc"
)

type server struct {
}

func (*server) GetSum(ctx context.Context, req *calpb.SumRequest) (*calpb.SumResponse, error) {
	sum := req.GetNum1() + req.GetNum2()
	result := &calpb.SumResponse{
		Sum: sum,
	}
	return result, nil
}

func (*server) GetAvg(stream calpb.Calculator_GetAvgServer) error {
	avg := float32(0.0)
	i := 0
	for {
		val, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calpb.AverageResponse{
				AvgRes: avg / float32(i),
			})
		}
		if err != nil {
			log.Fatal("Failed to process request: ", err)
		}

		log.Printf("Recieve %v from client.", val.GetAvgNums())
		avg += float32(val.AvgNums)
		i++
	}
}

func (*server) GetPrimeNum(req *calpb.PrimeRequest, serverStream calpb.Calculator_GetPrimeNumServer) error {
	k := int64(2)
	n := req.GetPrimeNum()
	for n > 1 {
		if n%k == 0 {
			n = n / k
			serverStream.Send(&calpb.PrimeResponse{
				PrimeResult: k,
			})
		} else {
			k = k + 1
		}
	}
	return nil
}

func main() {
	log.Println("-----Starting server!-----")
	conn, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("Failed to listen: ", err)
	}

	s := grpc.NewServer()
	calpb.RegisterCalculatorServer(s, &server{})
	if err := s.Serve(conn); err != nil {
		log.Fatal("Failed to start server!")
	}
}
