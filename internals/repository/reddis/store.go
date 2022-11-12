package reddis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"sync"
)

type Store struct {
	client *redis.Client
	ctx    *context.Context
	Keys   []string
	mu     sync.Mutex
}

func NewStore(rdb *redis.Client) *Store {
	return &Store{client: rdb}
}
func (s *Store) SaveMsg(key string, value string) error {

	s.mu.Lock()
	s.Keys = append(s.Keys, key)
	s.mu.Unlock()

	err := s.client.Set(key, value, 0).Err()
	if err != nil {
		return fmt.Errorf("cannot save the msg in store: %w", err)
	}
	return nil
}

func (s *Store) GetAllMsg() (string, error) {
	msg := ""

	s.mu.Lock()
	for _, key := range s.Keys {
		val, err := s.client.Get(key).Result()
		if err != nil {
			return msg, err
		}
		msg += val + "\n"
	}
	s.mu.Unlock()
	return msg, nil
}
