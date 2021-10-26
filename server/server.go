package main

import (
	Chat "Chitty_Chat/Chat"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type ChatServer struct {
	Chat.UnimplementedChittyChatServiceServer
}

var messageList []*Chat.FromClient

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

	// receive request <<< client
	go receiveFromStream(ccsi)

	//stream >>> client
	errch := make(chan error)
	go sendToStream(ccsi, errch)

	return <-errch
}

// receive from stream
func receiveFromStream(ccsi_ Chat.ChittyChatService_GetServerStreamServer) {
	for {
		req, err := ccsi_.Recv()
		if err != nil {
			log.Printf("Error reciving request from client :: %v", err)
			break

		} else {
			messageList = append(messageList, req)
			log.Printf("%v", len(messageList))
		}
	}
}

//send to stream
func sendToStream(ccsi_ Chat.ChittyChatService_GetServerStreamServer, errch_ chan error) {
	for {

		for {

			time.Sleep(500 * time.Millisecond)

			if len(messageList) == 0 {
				break
			}

			message := messageList[0]

			err := ccsi_.Send(&Chat.FromServer{Name: message.Name, Body: message.Body})

			if err != nil {
				errch_ <- err
			}

			if len(messageList) >= 2 {
				messageList = messageList[1:]
			} else {
				messageList = []*Chat.FromClient{}
			}
		}

		time.Sleep(1 * time.Second)
	}
}
