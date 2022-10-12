package lru

import "container/list"

type Cache struct {
	maxBytes int64
	nBytes   int64
	ll       *list.List
	cache    map[string]*list.Element

	// 记录被移除时的回调函数
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

//为了通用性，我们允许值是实现了 Value 接口的任意类型，
//该接口只包含了一个方法 Len() int，
//用于返回值所占用的内存大小。
type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}

	// 默认为false
	return
}

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())

		// callback
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nBytes += int64(len(key)) + int64(value.Len())
	}

	// 如果最大可用内存不为0  一直淘汰最老元素 只到已用内存小于最大要求
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}

}

// Len is the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
