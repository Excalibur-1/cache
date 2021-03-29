package cache

import (
	"container/list"
	"context"
	"fmt"
	"sync"
	"time"
)

func init() {
	fmt.Println("Register Mem to Cache Engine ver:1.0.0")
	Register(Mem, &defaultProvider{})
}

type defaultProvider struct {
	mutex       sync.RWMutex
	cacheList   *list.List
	cache       map[string]*list.Element
	maxItemSize int
}

type entry struct {
	key, value string
}

func (d *defaultProvider) New(ctx context.Context, cfg Config) (prov Provider, err error) {
	fmt.Println("Loading Memory Cache Engine ver:1.0.0")
	prov = &defaultProvider{
		cacheList:   list.New(),
		cache:       make(map[string]*list.Element),
		maxItemSize: cfg.MaxItemSize,
	}
	return
}

func (d *defaultProvider) Del(ctx context.Context, key ...string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if d.cache == nil {
		return
	}
	for _, v := range key {
		if ele, ok := d.cache[v]; ok {
			d.cacheList.Remove(ele)
			delete(d.cache, ele.Value.(*entry).key)
			return
		}
	}
}

func (d *defaultProvider) Exists(ctx context.Context, key string) bool {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	_, hit := d.cache[key]
	return hit
}

func (d *defaultProvider) Get(ctx context.Context, key string) (val string) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if ele, hit := d.cache[key]; hit {
		d.cacheList.MoveToFront(ele)
		val = ele.Value.(*entry).value
	}
	return
}

func (d *defaultProvider) Set(ctx context.Context, key, value string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.cache == nil {
		d.cache = make(map[string]*list.Element)
		d.cacheList = list.New()
	}

	if ele, ok := d.cache[key]; ok {
		d.cacheList.MoveToFront(ele)
		ele.Value.(*entry).value = value
		return
	}

	ele := d.cacheList.PushFront(&entry{key: key, value: value})
	d.cache[key] = ele
	if d.maxItemSize != 0 && d.cacheList.Len() > d.maxItemSize {
		d.removeOldest()
	}
}

func (d *defaultProvider) SetExpires(ctx context.Context, key, value string, expires time.Duration) {
	panic("implement me")
}

func (d *defaultProvider) HDel(ctx context.Context, key string, fields ...string) {
	panic("implement me")
}

func (d *defaultProvider) HExists(ctx context.Context, key, field string) bool {
	panic("implement me")
}

func (d *defaultProvider) HSet(ctx context.Context, key, field, value string) {
	d.Set(ctx, key+field, value)
}

func (d *defaultProvider) HGet(ctx context.Context, key, field string) string {
	return d.Get(ctx, key+field)
}

func (d *defaultProvider) HGetAll(ctx context.Context, key string) map[string]string {
	panic("implement me")
}

func (d *defaultProvider) Val(ctx context.Context, script string, keys []string, args ...interface{}) {
	panic("implement me")
}

func (d *defaultProvider) TTL(ctx context.Context, key string) int {
	panic("implement me")
}

func (d *defaultProvider) Close() {
	panic("implement me")
}

// RemoveOldest remove the oldest key
func (d *defaultProvider) removeOldest() {
	if d.cache == nil {
		return
	}
	ele := d.cacheList.Back()
	if ele != nil {
		d.cacheList.Remove(ele)
		key := ele.Value.(*entry).key
		delete(d.cache, key)
	}
}
