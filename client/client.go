package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"Chitty_Chat/ChittyChatService"

	"google.golang.org/grpc"
)

type clientHandle struct {
	stream     ChittyChatService.ChittyChatService_ChittyChatServiceClient
	clientName string
}

func main() {

	const serverID = "localhost:8008"

	log.Println("Connecting : " + serverID)
	conn, err := grpc.Dial(serverID, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Failed to connect gRPC server :: %v", err)
	}
	defer conn.Close()

	client := ChittyChatService.NewChittyChatServiceClient(conn)

	_stream, err := client.ChittyChatService(context.Background())
	if err != nil {
		log.Fatalf("Failed to get response from gRPC server :: %v", err)
	}

	ch := clientHandle{stream: _stream}
	ch.clientConfig()
	go ch.sendMessage()
	go ch.receiveMessage()

	// block main
	bl := make(chan bool)
	<-bl

}

// Assign name
func (ch *clientHandle) clientConfig() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Your Name : ")
	msg, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read from console :: %v", err)

	}
	ch.clientName = strings.TrimRight(msg, "\r\n")

}

func (ch *clientHandle) sendMessage() {

	for {
		reader := bufio.NewReader(os.Stdin)
		clientMessage, err := reader.ReadString('\n')
		clientMessage = strings.TrimRight(clientMessage, "\r\n")
		if err != nil {
			log.Printf("Failed to read from console :: %v", err)
			continue
		}

		clientMessageBox := &ChittyChatService.FromClient{
			Name: ch.clientName,
			Body: clientMessage,
		}

		err = ch.stream.Send(clientMessageBox)

		if err != nil {
			log.Printf("Error while sending to server :: %v", err)
		}

	}

}

func (ch *clientHandle) receiveMessage() {

	for {
		resp, err := ch.stream.Recv()
		if err != nil {
			log.Fatalf("can not receive1 %v", err)
		}
		log.Printf("%s : %s", resp.Name, resp.Body)
	}
}
