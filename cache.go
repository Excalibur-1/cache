package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var (
	providersMu sync.RWMutex
	providers   = make(map[RegisterName]Provider)
)

type RegisterName string

const (
	Redis RegisterName = "redis"
	Mem   RegisterName = "mem"
)

type Config struct {
	Provider    RegisterName `json:"provider"`      // 缓存实例，多种可选：redis、memcached等
	Host        string       `json:"host"`          // 地址
	Port        string       `json:"port"`          // 端口
	Password    string       `json:"password"`      // 密码
	MaxItemSize int          `json:"max_item_size"` // 默认缓存实现方案的可缓存元素上限
}

func Register(name RegisterName, provider Provider) {
	providersMu.Lock()
	defer providersMu.Unlock()
	if provider == nil {
		panic("cache: Register provider is nil")
	}
	// 判断是否已有相同的provider实现，有则报错退出
	if _, dup := providers[name]; dup {
		panic("cache: Register called twice for provider " + name)
	}
	providers[name] = provider
}

// Engine 缓存引擎定义，通过缓存引擎获取缓存客户端并进行数据缓存操作。系统根据业务类型划分了几个特定的缓存节点类型，每个类型的
// 缓存节点可以分别指定自己的缓存实现方式，通过配置中心的配置示例如下：
// 缓存节点服务器的类型，不同节点类型缓存的数据及其目的有所差异，业务系统要根据实际情况进行选择处理。
func Engine(ctx context.Context, conf Config) (prov Provider) {
	fmt.Println("Loading Cache Engine ver:1.0.0")
	provider := conf.Provider
	if provider == "" {
		panic("没有定义缓存实现信息!")
	}
	prov, ok := providers[provider]
	if !ok {
		panic("加载缓存实现出错!")
	}
	prov, err := prov.New(ctx, conf)
	if err != nil {
		panic("ping 缓存服务失败!")
	}
	return
}

// Provider 分布式的缓存操作接口。
type Provider interface {
	// New 创建新的provider实例
	New(ctx context.Context, cfg Config) (prov Provider, err error)
	// Del 从缓存中删除指定key的缓存数据。
	Del(ctx context.Context, key ...string)
	// Exists 判断缓存中是否存在指定的key
	Exists(ctx context.Context, key string) bool
	// Get 根据给定的key从分布式缓存中读取数据并返回，如果不存在或已过期则返回空
	Get(ctx context.Context, key string) string
	// Set 使用指定的key将对象存入分布式缓存中，并使用缓存的默认过期设置，注意，存入的对象必须是可序列化的。
	Set(ctx context.Context, key, value string)
	// SetExpires 使用指定的key将对象存入分部式缓存中，并指定过期时间，注意，存入的对象必须是可序列化的
	SetExpires(ctx context.Context, key, value string, expires time.Duration)
	// HDel 将指定key的map数据中的某个字段删除。
	HDel(ctx context.Context, key string, fields ...string)
	// HExists 判断缓存中指定key的map是否存在指定的字段，如果key或字段不存在则返回false。
	HExists(ctx context.Context, key, field string) bool
	// HSet 将指定key的map数据的某个字段设置为给定的值。
	HSet(ctx context.Context, key, field, value string)
	// HGet 获取指定key的map数据某个字段的值，如果不存在则返回空
	HGet(ctx context.Context, key, field string) string
	// HGetAll 获取指定key的map对象，如果不存在则返回nil
	HGetAll(ctx context.Context, key string) map[string]string
	// Val 对指定的key结果集执行指定的脚本
	Val(ctx context.Context, script string, keys []string, args ...interface{})
	// TTL 返回key的剩余存活时间
	TTL(ctx context.Context, key string) int
	// Close 关闭客户端
	Close()
}
