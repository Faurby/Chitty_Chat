package main

import (
	Chat "Chitty_Chat/Chat"
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

var lamport int32

type ChatServer struct {
	Chat.UnimplementedChittyChatServiceServer
}

var users = make(map[int32]Chat.ChittyChatService_JoinChatServer)

func main() {
	lamport = 0

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
	incomingLamport := user.Lamport

	lamport = max(lamport, incomingLamport)
	lamport++

	users[user.Id] = ccsi

	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic: %v", err)
			os.Exit(1)
		}
	}()

	body := fmt.Sprintf(user.Name + " has joined the chat")

	Broadcast(&Chat.FromClient{Name: "ServerMessage", Body: body, Lamport: lamport})

	// block function
	bl := make(chan bool)
	<-bl

	return nil
}

func (s *ChatServer) LeaveChat(ctx context.Context, user *Chat.User) (*Chat.Empty, error) {
	incomingLamport := user.Lamport
	lamport = max(lamport, incomingLamport)
	lamport++

	delete(users, user.Id)
	// updateTimestamp(int(user.Timestamp))
	Broadcast(&Chat.FromClient{Name: "ServerMessage", Body: user.Name + " has left the chat at Lamport time ", Lamport: lamport})
	return &Chat.Empty{}, nil
}

//ChatService
func (is *ChatServer) Publish(ctx context.Context, msg *Chat.FromClient) (*Chat.Empty, error) {

	incomingLamport := msg.Lamport
	lamport = max(lamport, incomingLamport)
	lamport++

	msg.Lamport = lamport

	Broadcast(msg)
	return &Chat.Empty{}, nil
}

func Broadcast(msg *Chat.FromClient) {
	name := msg.Name
	body := msg.Body
	time := msg.Lamport

	log.Printf("%s : %s %d", name, body, time)

	for key, value := range users {
		err := value.Send(&Chat.FromServer{Name: name, Body: body, Lamport: time})
		if err != nil {
			log.Println("Failed to broadcast to "+string(key)+": ", err)
		}
	}
}

func max(x, y int32) int32 {
	if x > y {
		return x
	} else {
		return y
	}
}
