package main

import (
	Chat "Chitty_Chat/Chat"
	"context"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

type ChatServer struct {
	Chat.UnimplementedChittyChatServiceServer
}

var users = make(map[int32]Chat.ChittyChatService_JoinChatServer)

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

func (c *ChatServer) JoinChat(user *Chat.User, ccsi Chat.ChittyChatService_JoinChatServer) error {
	users[user.Id] = ccsi

	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic: %v", err)
			os.Exit(1)
		}
	}()

	Broadcast(&Chat.FromClient{Name: "ServerMessage", Body: user.Name + " has joined the chat"})

	// block function
	bl := make(chan bool)
	<-bl

	return nil
}

func (s *ChatServer) LeaveChat(ctx context.Context, user *Chat.User) (*Chat.Empty, error) {
	delete(users, user.Id)
	// updateTimestamp(int(user.Timestamp))
	Broadcast(&Chat.FromClient{Name: "ServerMessage", Body: user.Name + " has left the chat"})
	return &Chat.Empty{}, nil
}

//ChatService
func (is *ChatServer) Publish(ctx context.Context, msg *Chat.FromClient) (*Chat.Empty, error) {
	Broadcast(msg)
	return &Chat.Empty{}, nil
}

func Broadcast(msg *Chat.FromClient) {
	name := msg.Name
	body := msg.Body

	log.Println(name + ": " + body)

	for key, value := range users {
		err := value.Send(&Chat.FromServer{Name: name, Body: body})
		if err != nil {
			log.Println("Failed to broadcast to "+string(key)+": ", err)
		}
	}
}
