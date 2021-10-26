package main

import (
	"log"
	"net"

	"Chitty_Chat/ChittyChatService"

	"google.golang.org/grpc"
)

func main() {
	listen, err := net.Listen("tcp", ":8008")
	if err != nil {
		log.Fatalf("Failed to listen on port 8008: %v", err)
	}
	log.Println(":8008 is listening")

	// Creates empty gRPC server
	grpcServer := grpc.NewServer()

	// Creates instance of our ChittyChatServer struct and binds it with our empty gRPC server.
	ccs := ChittyChatService.ChatServer{}
	ChittyChatService.RegisterChittyChatServiceServer(grpcServer, &ccs)

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("Failed to start gRPC server :: %v", err)
	}

}
