package server

import (
	"ZakirAvrora/ChatRoom/internals/models"
	"ZakirAvrora/ChatRoom/internals/repository/reddis"
	"github.com/go-redis/redis"
	"sync"
)

type Server struct {
	Rooms map[string]*models.ChatRoom
	rdb   *redis.Client
	Mu    sync.RWMutex
}

func NewServer(rdb *redis.Client) *Server {
	return &Server{
		Rooms: make(map[string]*models.ChatRoom),
		rdb:   rdb,
	}
}

func (s *Server) CreateNewRoom(roomName string, cap int) *models.ChatRoom {
	s.Mu.Lock()
	store := reddis.NewStore(roomName+"Chat", s.rdb)
	room := models.NewChatRoom(roomName, cap, store)
	s.Rooms[roomName] = room
	s.Mu.Unlock()

	go room.RunChatRoom()

	return room
}
