package reddis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
)

type Store struct {
	client *redis.Client
	ctx    *context.Context
}

func NewStore(rdb *redis.Client) *Store {
	return &Store{client: rdb}
}
func (s *Store) SaveMsg(msg string) error {
	if err := s.client.RPush("chat_messages", msg).Err(); err != nil {
		return fmt.Errorf("cannot save the msg in store: %w", err)
	}
	return nil
}

func (s *Store) GetAllMsg() ([]string, error) {
	msg, err := s.client.LRange("chat_messages", 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("cannot get the previoise msg in store: %w", err)
	}
	return msg, nil
}
