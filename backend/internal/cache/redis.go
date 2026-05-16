package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	DefaultTTL = 7 * 24 * time.Hour // 7 days
	KeyPrefix  = "blog:"
)

type RedisCache struct {
	client *redis.Client
}

func New(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

// Key builders

func PostKey(slug string) string {
	return fmt.Sprintf("%spost:%s", KeyPrefix, slug)
}

func PostListKey(tag string, page, limit int) string {
	return fmt.Sprintf("%sposts:list:%s:%d:%d", KeyPrefix, tag, page, limit)
}

func TagsKey() string {
	return KeyPrefix + "tags:all"
}

// Get retrieves a cached value by key. Returns empty string on miss.
func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

// Set stores a value with the default TTL.
func (c *RedisCache) Set(ctx context.Context, key, value string) error {
	return c.client.Set(ctx, key, value, DefaultTTL).Err()
}

// FlushBlog deletes all blog:* keys. Called after sync.
func (c *RedisCache) FlushBlog(ctx context.Context) error {
	iter := c.client.Scan(ctx, 0, KeyPrefix+"*", 100).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return err
	}
	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}
	return nil
}
