package singleflight

import "sync"

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	// lazy load
	if g.m == nil {
		g.m = make(map[string]*call)
	}

	// 此前已经有请求
	if c, ok := g.m[key]; ok {
		// 释放掉的m的锁
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	// 第一次请求
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	// 请求处理函数
	c.val, c.err = fn()
	c.wg.Done()

	// 将m中对应的key删除掉
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
