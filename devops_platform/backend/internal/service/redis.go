package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisService struct {
	client *redis.Client
	mu     sync.RWMutex
}

var redisService *RedisService

func GetRedisService() *RedisService {
	if redisService == nil {
		redisService = &RedisService{}
	}
	return redisService
}

func (s *RedisService) Connect(host, port, password string, db int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("连接Redis失败: %w", err)
	}

	s.client = client
	return nil
}

func (s *RedisService) GetKeys(ctx context.Context, pattern string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.client == nil {
		return nil, fmt.Errorf("Redis未连接")
	}

	if pattern == "" {
		pattern = "*"
	}

	keys, err := s.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (s *RedisService) GetValue(ctx context.Context, key string) (map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.client == nil {
		return nil, fmt.Errorf("Redis未连接")
	}

	keyType, err := s.client.Type(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"key":  key,
		"type": keyType,
	}

	ttl, _ := s.client.TTL(ctx, key).Result()
	result["ttl"] = int64(ttl.Seconds())

	switch keyType {
	case "string":
		val, err := s.client.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		result["value"] = val

	case "list":
		val, err := s.client.LRange(ctx, key, 0, -1).Result()
		if err != nil {
			return nil, err
		}
		result["value"] = val

	case "set":
		val, err := s.client.SMembers(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		result["value"] = val

	case "hash":
		val, err := s.client.HGetAll(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		result["value"] = val

	case "zset":
		val, err := s.client.ZRangeWithScores(ctx, key, 0, -1).Result()
		if err != nil {
			return nil, err
		}
		result["value"] = val

	default:
		result["value"] = nil
	}

	return result, nil
}

func (s *RedisService) SetValue(ctx context.Context, key, value string, ttl int) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.client == nil {
		return fmt.Errorf("Redis未连接")
	}

	expiration := time.Duration(0)
	if ttl > 0 {
		expiration = time.Duration(ttl) * time.Second
	}

	return s.client.Set(ctx, key, value, expiration).Err()
}

func (s *RedisService) DeleteKey(ctx context.Context, key string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.client == nil {
		return fmt.Errorf("Redis未连接")
	}

	return s.client.Del(ctx, key).Err()
}

func (s *RedisService) GetInfo(ctx context.Context) (map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.client == nil {
		return nil, fmt.Errorf("Redis未连接")
	}

	info, err := s.client.Info(ctx).Result()
	if err != nil {
		return nil, err
	}

	dbSize, _ := s.client.DBSize(ctx).Result()

	return map[string]interface{}{
		"info":   info,
		"dbSize": dbSize,
	}, nil
}
