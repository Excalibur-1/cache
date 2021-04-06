package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Excalibur-1/configuration"
	"testing"
)

func TestReadConfiguration(t *testing.T) {
	mockEngine := configuration.MockEngine(map[string]string{
		"/myconf/base/cache/1000": "{\"provider\":\"redis\",\"host\":\"127.0.0.1\",\"port\":\"6379\",\"password\":\"\"}",
	})
	s, err := mockEngine.String("myconf", "base", "cache", "", "1000")
	if err != nil {
		panic(err)
	}
	var conf Config
	if err = json.Unmarshal([]byte(s), &conf); err != nil {
		panic(err)
	}
	ctx := context.Background()
	eng := Engine(ctx, conf)
	get := eng.Get(ctx, "test")
	fmt.Println(get)
}

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
