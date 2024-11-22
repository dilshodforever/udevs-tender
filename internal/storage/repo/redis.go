package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func (s *BidService) GetCache(ctx context.Context, key string) ([]byte, error) {
	// Redis-dan ma'lumot o'qish
	cachedData, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("cache miss for key: %s", key)
		}
		return nil, fmt.Errorf("failed to get cache for key %s: %v", key, err)
	}
	return []byte(cachedData), nil
}

func (s *BidService) SetCache(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value for cache key %s: %v", key, err)
	}

	if err := s.redis.Set(ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set cache for key %s: %v", key, err)
	}

	return nil
}
