package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCacheService struct {
	ctx         context.Context
	redisClient *redis.Client
}

func NewRedisCacheService(redisClient *redis.Client) RedisCacheService {
	return &redisCacheService{
		ctx:         context.Background(),
		redisClient: redisClient,
	}
}

func (cs *redisCacheService) Get(key string, dest any) error {
	data, err := cs.redisClient.Get(cs.ctx, key).Result()
	if err == redis.Nil {
		return err
	}
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

func (cs *redisCacheService) Set(key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return cs.redisClient.Set(cs.ctx, key, data, ttl).Err()
}

func (cs *redisCacheService) Clear(pattern string) error {
	cursor := uint64(0)
	for {
		keys, nextCursor, err := cs.redisClient.Scan(cs.ctx, cursor, pattern, 2).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			cs.redisClient.Del(cs.ctx, keys...)
		}

		cursor = nextCursor

		if cursor == 0 {
			break
		}
	}

	return nil
}

func (cs *redisCacheService) Exists(key string) (bool, error) {
	count, err := cs.redisClient.Exists(cs.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (cs *redisCacheService) Del(key string) error {
	return cs.redisClient.Del(cs.ctx, key).Err()
}
