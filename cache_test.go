package cache

import (
	"context"
	"fmt"
	"testing"
)

func TestSingleRedisProvider_Set(t *testing.T) {
	ctx := context.Background()
	eng := Engine(ctx, Config{
		Provider: Redis,
		Host:     "localhost",
		Port:     "6379",
	})
	eng.Set(ctx, "test", "test test")
	get := eng.Get(ctx, "test")
	fmt.Println(get)
}

func TestDefaultProvider_Set(t *testing.T) {
	ctx := context.Background()
	eng := Engine(ctx, Config{
		Provider:    Mem,
		MaxItemSize: 1000,
	})
	eng.Set(ctx, "test", "test test")
	get := eng.Get(ctx, "test")
	fmt.Println(get)
}
