package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"

	"github.com/grpc-go-course/greet/greetpb"
)

type server struct {
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	firstName := req.GetGreeting().GetFirstName()

	result := "Hello " + firstName

	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.Greet_GreetManyTimesServer) error {
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		resStr := "Hello " + firstName + " number: " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: resStr,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(req greetpb.Greet_LongGreetServer) error {
	resStr := ""
	for {
		msg, err := req.Recv()
		if err == io.EOF {
			res := &greetpb.LongGreetResponse{
				Result: resStr,
			}
			req.SendAndClose(res)
		}
		if err != nil {
			log.Fatal("Error when proccess request first name: " + msg.GetGreeting().GetFirstName())
		}
		resStr += "Hello " + msg.GetGreeting().GetFirstName() + " ! "
	}
}

func main() {
	fmt.Println("Starting server !!!!!")

	conn, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatal("Failed to listen: ", err)
	}

	s := grpc.NewServer()

	greetpb.RegisterGreetServer(s, &server{})

	if err := s.Serve(conn); err != nil {
		log.Fatal("Failed to start grpc server", err)
	}

}
