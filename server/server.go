package server

import (
	"ZakirAvrora/ChatRoom/internals/models"
)

type Server struct {
	Rooms map[string]*models.ChatRoom
}

func NewServer() *Server {
	return &Server{
		Rooms: make(map[string]*models.ChatRoom),
	}
}
