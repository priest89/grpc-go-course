package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting client !!!!!")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatal("Failed to connect server: ", err)
	}

	defer cc.Close()

	greetClient := greetpb.NewGreetClient(cc)

	// doUnary(greetClient)
	// doGreetManyTimes(greetClient)
	doLongGreet(greetClient)
}

func doLongGreet(greetClient greetpb.GreetClient) {
	reqs := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Priest1",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Vu89",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Thang12",
			},
		},
	}

	stream, err := greetClient.LongGreet(context.Background())

	for _, req := range reqs {
		log.Printf("Send req %v", req.GetGreeting().GetFirstName())
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatal("Failed when get response from server ", err)
	}

	fmt.Println(res.GetResult())
}

func doGreetManyTimes(greetClient greetpb.GreetClient) {
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Priest",
			LastName:  "Vu",
		},
	}

	resStream, err := greetClient.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatal("Failed to get response from server: ", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error while reading stream: ", err)
		}
		fmt.Println("Response from GreetManyTimes: ", msg.GetResult())
	}
}

func doUnary(greetClient greetpb.GreetClient) {

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Thang",
			LastName:  "Vu",
		},
	}

	res, err := greetClient.Greet(context.Background(), req)

	if err != nil {
		log.Fatal("Failed to get data from server", err)
	}

	fmt.Println(res.GetResult())
}
