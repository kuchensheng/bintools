package cacher

import (
	"github.com/kuchensheng/bintools/cache/store"
	"time"
)

//CacheInterface return the interface for all caches
type CacheInterface interface {
	Get(key string) (any, bool)
	GetWithExpiration(key string) (any, time.Time, bool)
	Set(key string, value any, options ...store.Options) error
	Expiration(key string, options ...store.Options) error
	GetTTL(key string) (int64, error)
	Delete(key string) error
	Clear() error
}
