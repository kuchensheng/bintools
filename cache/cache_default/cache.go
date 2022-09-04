package cache_default

import (
	"io"
	"runtime"
	"sync"
	"time"
)

type DefaultCache *cache

type Item struct {
	//Data
	Data any
	//Ttl time.UnixNano.Expiration of Item
	Ttl int64
}

type janitor struct {
	Interval time.Duration
	stop     chan bool
}

type cache struct {
	defaultExpiration time.Duration
	items             map[string]Item
	mu                sync.RWMutex
	//如果被lru算法驱逐，执行该方法
	onEvicted func(string, any)
	j         *janitor
}

//Cap 计算当前缓存的长度，加锁
func (c *cache) Cap() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	ci := c.items
	return len(ci)
}

func (c *cache) getUnixNano(expiration time.Duration) int64 {
	de := c.defaultExpiration
	if expiration > -1 {
		de = expiration
	}
	e := int64(-1)
	if de > -1 {
		e = time.Now().Add(de).UnixNano()
	}
	return e
}

//Expired 判断当前元素是否过期
func (item *Item) Expired() bool {
	if item.Ttl == -1 {
		return false
	}
	return time.Now().UnixNano() > item.Ttl
}

func New() *cache {
	return NewWithExpiration(-1)
}

func NewWithExpiration(expiration time.Duration) *cache {
	return NewWithExpirationAndCleanupInterval(expiration, 0)
}

func NewWithExpirationAndCleanupInterval(defaultExpiration, cleanupInterval time.Duration) *cache {
	if defaultExpiration <= 0 {
		defaultExpiration = -1
	}
	ch := make(chan bool)
	c := &cache{
		defaultExpiration: defaultExpiration,
		items:             make(map[string]Item),
		j: &janitor{
			Interval: 500 * time.Millisecond,
			stop:     ch,
		},
	}
	c.OnEvicted(c.save)
	//启动自动清理器
	go func() {
		c.runCleanup(cleanupInterval)
		runtime.SetFinalizer(c, stopJanitor)
	}()
	return c
}

func (c *cache) runCleanup(cleanupInterval time.Duration) {
	if cleanupInterval <= 0 {
		cleanupInterval = 500 * time.Millisecond
	}
	ticker := time.NewTicker(cleanupInterval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-c.j.stop:
			ticker.Stop()
			return
		}
	}
}

func stopJanitor(c *cache) {
	c.j.stop <- true
}

//DeleteExpired 删除过期的key
func (c *cache) DeleteExpired() {
	l := len(c.items)
	if l < 1 {
		//无需处理
		return
	}
	//fixme 这里将会引起内存增大
	cloneMap := make(map[string]Item, c.Cap())
	for k, v := range c.items {
		cloneMap[k] = v
	}
	//这里的长度，需要再次判断c.items
	ch := make(chan int8, len(c.items))
	for key, item := range cloneMap {
		//开启多协程，快速处理
		go func(k string, i Item) {
			c.mu.Lock()
			defer c.mu.Unlock()
			if i.Expired() {
				delete(c.items, k)
			}
			ch <- int8(1)
		}(key, item)
	}
}

func (c *cache) save(key string, value any) {
	//todo 存储起来
	//gob.NewDecoder()
}

func (c *cache) OnEvicted(f func(key string, value any)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onEvicted = f
}

func (c *cache) Save(w io.Writer) error {
	return nil
}
