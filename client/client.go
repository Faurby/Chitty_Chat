package main

import (
	Chat "Chitty_Chat/Chat"
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

type clientHandle struct {
	stream     Chat.ChittyChatService_JoinChatClient
	Id         int32
	clientName string
}

func main() {

	const serverID = "localhost:8007"

	log.Println("Connecting : " + serverID)
	conn, err := grpc.Dial(serverID, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Failed to connect gRPC server :: %v", err)
	}
	defer conn.Close()

	client := Chat.NewChittyChatServiceClient(conn)

	ch := clientHandle{}
	ch.clientConfig()

	rand.Seed(time.Now().UnixNano())
	ch.Id = rand.Int31()

	var user = &Chat.User{
		Id:   ch.Id,
		Name: ch.clientName,
	}

	_stream, err := client.JoinChat(context.Background(), user)
	if err != nil {
		log.Fatalf("Failed to get response from gRPC server :: %v", err)
	}

	ch.stream = _stream

	SetupCloseHandler(ch, client)

	go ch.sendMessage(client)
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

func (ch *clientHandle) sendMessage(client Chat.ChittyChatServiceClient) {

	for {
		reader := bufio.NewReader(os.Stdin)
		clientMessage, err := reader.ReadString('\n')
		clientMessage = strings.TrimRight(clientMessage, "\r\n")
		if err != nil {
			log.Printf("Failed to read from console : %v", err)
			continue
		}

		msg := &Chat.FromClient{
			Name: ch.clientName,
			Body: clientMessage,
		}

		_, err = client.Publish(context.Background(), msg)
		if err != nil {
			log.Printf("Error while sending to server :: %v", err)
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func (ch *clientHandle) receiveMessage() {

	for {
		resp, err := ch.stream.Recv()
		if err != nil {
			log.Fatalf("can not receive %v", err)
		}
		log.Printf("%s : %s", resp.Name, resp.Body)
	}
}

// SetupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS. We then handle this by calling
// our clean up procedure and exiting the program.
func SetupCloseHandler(ch clientHandle, client Chat.ChittyChatServiceClient) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		client.LeaveChat(context.Background(), &Chat.User{Id: ch.Id, Name: ch.clientName})
		os.Exit(0)
	}()
}
