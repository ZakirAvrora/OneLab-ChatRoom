package reddis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
)

type Store struct {
	name   string
	client *redis.Client
	ctx    *context.Context
}

func NewStore(name string, rdb *redis.Client) *Store {
	return &Store{name: name,
		client: rdb}
}
func (s *Store) SaveMsg(msg string) error {
	if err := s.client.RPush(s.name, msg).Err(); err != nil {
		return fmt.Errorf("cannot save the msg in store: %w", err)
	}
	return nil
}

func (s *Store) GetAllMsg() ([]string, error) {
	msg, err := s.client.LRange(s.name, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("cannot get the previose msg in store: %w", err)
	}
	return msg, nil
}
