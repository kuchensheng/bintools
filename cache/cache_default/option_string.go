package cache_default

import (
	"errors"
	"github.com/kuchensheng/bintools/cache/store"
	"time"
)

func (c *cache) Get(key string) (any, bool) {
	if item, found := c.get(key); !found {
		return nil, false
	} else if item.Expired() {
		//找到了过期的key,返回未找到信息
		return nil, false
	} else {
		return item.Data, true
	}
}

func (c *cache) get(key string) (Item, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if item, found := c.items[key]; !found {
		return Item{}, false
	} else {
		return item, true
	}
}

func (c *cache) GetWithExpiration(key string) (any, time.Time, bool) {
	if item, found := c.get(key); !found {
		return nil, time.Time{}, false
	} else if item.Expired() {
		return nil, time.Time{}, false
	} else {
		return item.Data, time.Unix(0, item.Ttl), true
	}
}

func (c *cache) Set(key string, value any, options ...store.Options) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	expiration := c.getUnixNano(-1)
	if len(options) > 0 {
		expiration = c.getUnixNano(options[0].Expiration)
	}
	c.items[key] = Item{
		Data: value,
		Ttl:  expiration,
	}
	return nil
}

func (c *cache) Expiration(key string, options ...store.Options) error {
	if len(options) < 1 {
		return errors.New("options is required")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if item, found := c.items[key]; found {
		//重新设置超时时间
		item.Ttl = c.getUnixNano(options[0].Expiration)
		return nil
	}
	return errors.New("not found")
}
func (c *cache) GetTTL(key string) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if item, found := c.items[key]; found {
		return item.Ttl, nil
	}
	return 0, errors.New("not found")
}
func (c *cache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
	return nil
}
func (c *cache) Clear() error {
	//清楚所有key信息
	c.mu.Lock()
	c.mu.Unlock()
	//10w以内的delete，速度与之媲美。这是因为为delete会转换成mapclear
	c.items = make(map[string]Item)
	return nil
}
