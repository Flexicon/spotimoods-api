package cache

import (
	"context"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/flexicon/spotimoods-go/internal"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// adapter for a github.com/go-redis/cache instance
type adapter struct {
	cache *cache.Cache
}

// NewCache builder
func NewCache() (internal.Cache, error) {
	r := redis.NewClient(&redis.Options{
		Addr: viper.GetString("redis.url"),
	})
	// Validate redis connection by pinging it, otherwise return err
	if err := r.Ping(context.Background()).Err(); err != nil {
		return nil, errors.Wrap(err, "cache setup failed")
	}

	c := cache.New(&cache.Options{
		Redis:      r,
		LocalCache: fastcache.New(100 << 20), // 100 MB
	})

	return &adapter{
		cache: c,
	}, nil
}

func (a *adapter) Set(item *internal.CacheItem) error {
	return a.cache.Set(a.build(item))
}

func (a *adapter) Get(key string, value interface{}) error {
	return a.cache.Get(context.TODO(), key, &value)
}

func (a *adapter) Delete(key string) error {
	return a.cache.Delete(context.TODO(), key)
}

func (a *adapter) Exists(key string) bool {
	return a.cache.Exists(context.TODO(), key)
}

func (a *adapter) Once(item *internal.CacheItem) error {
	return a.cache.Once(a.build(item))
}

func (a *adapter) build(item *internal.CacheItem) *cache.Item {
	i := &cache.Item{
		Key:            item.Key,
		Value:          item.Value,
		TTL:            item.TTL,
		IfExists:       item.IfExists,
		IfNotExists:    item.IfNotExists,
		SkipLocalCache: item.SkipLocalCache,
	}

	if item.Do != nil {
		i.Do = func(doItem *cache.Item) (interface{}, error) {
			return item.Do(a.buildInternal(doItem))
		}
	}

	return i
}

func (a *adapter) buildInternal(item *cache.Item) *internal.CacheItem {
	return &internal.CacheItem{
		Key:            item.Key,
		Value:          item.Value,
		TTL:            item.TTL,
		IfExists:       item.IfExists,
		IfNotExists:    item.IfNotExists,
		SkipLocalCache: item.SkipLocalCache,
	}
}
