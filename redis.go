package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func init() {
	fmt.Println("Register Redis to Cache Engine ver:1.0.0")
	Register(Redis, &SingleRedisProvider{})
}

// SingleRedisProvider 基于Redis单实例的缓存实现(单机）
type SingleRedisProvider struct {
	client *redis.Client
}

func (r *SingleRedisProvider) New(ctx context.Context, cfg Config) (prov Provider, err error) {
	r.client = redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
	})
	pin, err := r.client.Ping(ctx).Result()
	if err != nil {
		return
	}
	fmt.Printf("Redis Pong:%s\n", pin)
	prov = r
	return
}

func (r SingleRedisProvider) Del(ctx context.Context, key ...string) {
	r.client.Del(ctx, key...)
}

func (r SingleRedisProvider) Exists(ctx context.Context, key string) bool {
	_, err := r.client.Exists(ctx, key).Result()
	return err == redis.Nil
}

func (r SingleRedisProvider) Get(ctx context.Context, key string) string {
	return r.client.Get(ctx, key).Val()
}

func (r SingleRedisProvider) Set(ctx context.Context, key, value string) {
	r.SetExpires(ctx, key, value, -1)
}

func (r SingleRedisProvider) SetExpires(ctx context.Context, key, value string, expires time.Duration) {
	r.client.Set(ctx, key, value, expires)
}

func (r SingleRedisProvider) HDel(ctx context.Context, key string, fields ...string) {
	r.client.HDel(ctx, key, fields...)
}

func (r SingleRedisProvider) HExists(ctx context.Context, key, field string) bool {
	v, _ := r.client.HExists(ctx, key, field).Result()
	return v
}

func (r SingleRedisProvider) HSet(ctx context.Context, key, field, value string) {
	r.client.HSet(ctx, key, field, value)
}

func (r SingleRedisProvider) HGet(ctx context.Context, key, field string) string {
	return r.client.HGet(ctx, key, field).Val()
}

func (r SingleRedisProvider) HGetAll(ctx context.Context, key string) map[string]string {
	return r.client.HGetAll(ctx, key).Val()
}

func (r SingleRedisProvider) Val(ctx context.Context, script string, keys []string, args ...interface{}) {
	r.client.Eval(ctx, script, keys, args...)
}

func (r SingleRedisProvider) TTL(ctx context.Context, key string) (ttl int) {
	v, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return
	}
	ttl = int(v.Seconds())
	return
}

func (r SingleRedisProvider) Close() {
	_ = r.client.Close()
}
