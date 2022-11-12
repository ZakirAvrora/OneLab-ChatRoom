package models

import (
	"ZakirAvrora/ChatRoom/internals/repository"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
)

type ChatRoom struct {
	Name         string
	Members      map[*Client]bool
	MaxSize      int
	HistoryStore repository.Repository

	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
}

func NewChatRoom(name string, maxSize int, repo repository.Repository) *ChatRoom {
	return &ChatRoom{
		Name:         name,
		MaxSize:      maxSize,
		HistoryStore: repo,
		Members:      make(map[*Client]bool),
		Register:     make(chan *Client),
		Unregister:   make(chan *Client),
		Broadcast:    make(chan Message),
	}
}

func (r *ChatRoom) RunChatRoom() {
	for {
		select {
		case client := <-r.Register:
			r.registerMember(client)
		case client := <-r.Unregister:
			r.unregisterMember(client)
		case message := <-r.Broadcast:
			//fmt.Println("Broadcasting message: ", message)
			r.broadcastMsg(message)
		}
	}
}

func (r *ChatRoom) registerMember(client *Client) {
	history, err := r.HistoryStore.GetAllMsg()
	fmt.Println(history)
	if err != nil {
		log.Fatalln(err)
	}
	w, err := client.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}
	w.Write([]byte(history))
	r.Members[client] = true
}

func (r *ChatRoom) unregisterMember(client *Client) {
	_, ok := r.Members[client]
	if ok {
		delete(r.Members, client)
	}
}

func (r *ChatRoom) broadcastMsg(msg Message) {
	myUuid := uuid.New()
	fmt.Println(myUuid, string(msg.Msg))
	if err := r.HistoryStore.SaveMsg(myUuid.String(), string(msg.Msg)); err != nil {
		log.Println(err)
	}

	for member := range r.Members {
		member.Msg <- msg
	}
}
