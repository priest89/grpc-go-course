#bin bash

protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.

protoc calculator/pb/calculator.proto --go_out=plugins=grpc:.