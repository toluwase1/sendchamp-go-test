package db

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sendchamp-go-test/models"
)

type ServerSentEvent interface {
	CreateServerSentEvent(request *models.Task) ([]byte, error)
}

var wsconn *websocket.Conn

func NewSseRepo(wscon *websocket.Conn) {
	log.Println(wscon)
	wsconn = wscon
}

func CreateServerSentEvent(msg *models.Task) ([]byte, error) {
	msgBody, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}
	err = wsconn.WriteMessage(websocket.TextMessage, msgBody)
	if err != nil {
		fmt.Printf("error sending message: %s\n", err.Error())
		return []byte{}, err
	}
	return msgBody, nil
}
