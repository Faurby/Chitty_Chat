package ChittyChatService

import (
	"log"
	"time"
)

type ChatServer struct {
	UnimplementedChittyChatServiceServer
}

var messageList []*FromClient

//ChatService
func (is *ChatServer) ChittyChatService(ccsi ChittyChatService_ChittyChatServiceServer) error {

	// receive request <<< client
	go receiveFromStream(ccsi)

	//stream >>> client
	errch := make(chan error)
	go sendToStream(ccsi, errch)

	return <-errch
}

// receive from stream
func receiveFromStream(ccsi_ ChittyChatService_ChittyChatServiceServer) {
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
func sendToStream(ccsi_ ChittyChatService_ChittyChatServiceServer, errch_ chan error) {
	// for {

	for {

		time.Sleep(500 * time.Millisecond)

		message := messageList[0]

		err := ccsi_.Send(&FromServer{Name: message.Name, Body: message.Body})

		if err != nil {
			errch_ <- err
		}

		if len(messageList) >= 2 {
			messageList = messageList[1:]
		} else {
			messageList = []*FromClient{}
		}
	}

	// 	time.Sleep(1 * time.Second)
	// }
}
