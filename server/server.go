package main

import (
	Chat "Chitty_Chat/Chat"
	"log"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
)

type ChatServer struct {
	Chat.UnimplementedChittyChatServiceServer
}

var userAmount int
var users = make(map[string]Chat.ChittyChatService_GetServerStreamServer)

func main() {
	listen, err := net.Listen("tcp", ":8007")
	if err != nil {
		log.Fatalf("Failed to listen on port 8007: %v", err)
	}
	log.Println(":8007 is listening")

	// Creates empty gRPC server
	grpcServer := grpc.NewServer()

	// Creates instance of our ChittyChatServer struct and binds it with our empty gRPC server.
	ccs := ChatServer{}
	Chat.RegisterChittyChatServiceServer(grpcServer, &ccs)

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("Failed to start gRPC server :: %v", err)
	}

}

//ChatService
func (is *ChatServer) GetServerStream(ccsi Chat.ChittyChatService_GetServerStreamServer) error {

	userID := userAmount + 1
	userAmount++

	users[strconv.Itoa(userID)] = ccsi

	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic: %v", err)
			os.Exit(1)
		}
	}()

	for {
		input, error := ccsi.Recv()
		if error != nil {
			log.Fatalln("Fatal error:", error)
			break
		}

		Broadcast(input)
	}
	return nil
}

func Broadcast(msg *Chat.FromClient) {
	name := msg.Name
	body := msg.Body

	log.Println(name + ": " + body)

	for key, value := range users {
		err := value.Send(&Chat.FromServer{Name: name, Body: body})
		if err != nil {
			log.Println("Failed to broadcast to "+key+": ", err)
		}
	}
}
